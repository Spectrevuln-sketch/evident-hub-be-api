package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"evidence-hub-be/src/core/constants"
	"fmt"
	"log"
	math "math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func FormToJson(form *multipart.Form, data any) error {
	formData := make(map[string]interface{})

	for fieldName, values := range form.Value {
		if len(values) == 1 {
			formData[fieldName] = values[0]
		} else {
			formData[fieldName] = values
		}
	}

	jsonData, jsonErr := json.Marshal(formData)
	if jsonErr != nil {
		return jsonErr
	}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}
	return nil
}

// func ProcessExcel(pathFile string, result interface{}) error {
// 	xFile, err := xlsx.OpenFile(pathFile)
// 	if err != nil {
// 		return err
// 	}
// 	structType := reflect.TypeOf(result).Elem().Elem() // Get the element type of the slice
// 	var fieldNames []string
// 	for i := 0; i < structType.NumField(); i++ {
// 		field := structType.Field(i)
// 		fieldName := field.Tag.Get("json")
// 		fieldNames = append(fieldNames, fieldName)
// 	}
// 	headers := make(map[string]int)
// 	for i, cell := range xFile.Sheets[0].Rows[0].Cells {
// 		headers[cell.String()] = i
// 	}
// 	resultSlice := reflect.ValueOf(result).Elem() // Get the value of the result slice
// 	for _, row := range xFile.Sheets[0].Rows[1:] {
// 		newStruct := reflect.New(structType).Interface()
// 		for _, field := range fieldNames {
// 			columnIndex, exists := headers[field]
// 			if !exists {
// 				continue
// 			}
// 			fieldValue := reflect.ValueOf(newStruct).Elem().FieldByName(field)
// 			if fieldValue.IsValid() {
// 				switch fieldValue.Kind() {
// 				case reflect.String:
// 					fieldValue.SetString(row.Cells[columnIndex].String())
// 				case reflect.Int:
// 					age, _ := row.Cells[columnIndex].Int()
// 					fieldValue.SetInt(int64(age))
// 				}
// 			}
// 		}
// 		resultSlice.Set(reflect.Append(resultSlice, reflect.ValueOf(newStruct).Elem()))
// 	}
// 	return nil
// }

// ==============================Romain Types===================================
func ToRoman(n int) string {
	symbols := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	values := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}

	result := ""

	for i := 0; n > 0; i++ {
		for values[i] <= n {
			result += symbols[i]
			n -= values[i]
		}
	}

	return result
}

func FilterByProvinceID(regencies []map[string]interface{}, provinceID string) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, regency := range regencies {
		if regency["province_id"] == provinceID {
			filtered = append(filtered, regency)
		}
	}
	return filtered
}

func Contains(slice []string, item string) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

func MustAtoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return value
}

func StructToMap(input interface{}) (map[string]interface{}, error) {
	// Marshal the struct to JSON
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into a map
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func ReadTemplate(filepath string) (*excelize.File, error) {
	log.Println("Reading template from:", filepath)
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func FillTemplate(template *excelize.File, payloads []map[string]interface{}) {
	for i, item := range payloads {
		row := i + 2 // Assuming the first row is the header
		template.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), item["CoaCode"])
		template.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), item["ProductName"])
		template.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), item["QTY"])
		template.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), item["UOM"])
		template.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), item["Description"])
		template.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), item["RequestedBy"])
		template.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), item["Purpose"])
		// Continue for other fields
	}
}

func GenerateFilePathAndSave(
	file *multipart.FileHeader,
	typePath string,
	ctx *gin.Context,
) string {

	ext := filepath.Ext(file.Filename)

	fullPath := filepath.Join(
		constants.UploadDir,
		typePath,
		GenerateRandomString(5)+ext,
	)

	log.Printf("FILE NAME      : %s", file.Filename)
	log.Printf("FILE SIZE      : %d", file.Size)
	log.Printf("FULL PATH      : %s", fullPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		log.Printf("MKDIR ERROR    : %v", err)
		return ""
	}

	if err := ctx.SaveUploadedFile(file, fullPath); err != nil {
		log.Printf("SAVE ERROR     : %v", err)
		return ""
	}

	if _, err := os.Stat(fullPath); err != nil {
		log.Printf("STAT ERROR     : %v", err)
		return ""
	}

	log.Printf("SAVE SUCCESS   : %s", fullPath)

	return fullPath
}

func GeneratePDF(template *excelize.File, outputPath string) error {
	log.Println("Generating PDF from template")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Example of converting the first cell value to PDF
	cell, err := template.GetCellValue("Sheet1", "A1")
	if err != nil {
		return err
	}
	pdf.Cell(40, 10, cell)

	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return err
	}
	return nil
}

func GenerateRandomNumber() int {
	source := math.NewSource(time.Now().UnixNano())
	r := math.New(source)

	if r.Intn(2) == 0 {
		return 100 + r.Intn(900)
	}
	return 1000 + r.Intn(9000)
}

// Fungsi untuk menghasilkan NPWPD dengan struktur
func GenerateNPWPD(kodeGolongan, kodeWilayahDistrik, kodeWilayahKelurahan int) string {
	source := math.NewSource(time.Now().UnixNano())
	r := math.New(source)

	nomorPokok := r.Intn(10000000) // 7 digit (Nomor Pokok)

	npwpd := fmt.Sprintf("%02d%07d%02d%02d", kodeGolongan, nomorPokok, kodeWilayahDistrik, kodeWilayahKelurahan)

	return npwpd
}

func GenerateUsername(fullname, tandaSelar string) string {
	re := regexp.MustCompile(`GT.\s*(\d+)`)
	match := re.FindStringSubmatch(tandaSelar)

	if len(match) == 0 {
		return ""
	}

	reUser := regexp.MustCompile(`^(?:DR|D|IR|MM|SH|S\.H)\.?[\s,]*|[\s,]*(?:DR|D|IR|MM|SH|S\.H)\.?$|[^\w\s-]`)

	username := reUser.ReplaceAllString(strings.ToLower(fullname), " ")
	username = strings.ReplaceAll(username, " ", "")
	return fmt.Sprintf("%s%s", username, match[1])
}

// func ExcelGenerate(generatedName string, templateFilePath string, payloads []map[string]interface{}, ctx *fiber.Ctx) (string, error) {
// 	// Open the Excel template file from the doc folder
// 	templateFile, err := excelize.OpenFile(templateFilePath)
// 	if err != nil {
// 		return "", err
// 	}

// 	FillTemplate(templateFile, payloads)
// 	log.Printf("DATA excel file %v", templateFile)
// 	pdfFilePath, err := generatePDF(templateFile)
// 	if err != nil {
// 		return "", ResponseError("99", err.Error(), ctx)
// 	}
// 	return pdfFilePath, nil
// }
//
