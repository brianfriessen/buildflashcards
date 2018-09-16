package main

import (
	"bufio"
	"fmt"
	"github.com/brianfriessen/bingsearch"
	"github.com/brianfriessen/forvosearch"
	"github.com/nfnt/resize"
	"image"
	_ "image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func makestring(runeArray []byte) []string {
	var returnSlice []string

	for i := 0; i < len(runeArray); i++ {
		returnSlice = append(returnSlice, string(runeArray[i]))
	}
	return returnSlice
}

/* The purpose of this program is to take a list of words in a text file (1 word per line) and create a html page
which has 4 pictures of the word along with the word and a pronunciation.
*/
func check(e error) {
	if e != nil {
		fmt.Println("Huston we have a problem..")
		log.Fatal(e)
		panic(e)
	}
}

func writeHTMLHeader(file *os.File) {

	file.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"UTF-8\">\n")
	file.WriteString("<link rel=\"stylesheet\" href=\"main.css\">\n")
	file.WriteString("<title>Flash Card Builder</title>\n<body>\n")
	file.WriteString("<img src=\"./images/skull.jpg\"></head>\n")

}

func writeHTMLSound(file *os.File, mp3FileName string) {

	file.WriteString("<div>\n")
	file.WriteString("<label for= \"" + mp3FileName + "\">\n")
	file.WriteString("	" + mp3FileName + "\n")
	file.WriteString("</label>\n")
	file.WriteString("<br/>\n")
	file.WriteString("<audio\n")
	file.WriteString("\tid=\"" + mp3FileName + "-sound\"\n")
	file.WriteString("\tcontrols\n")
	file.WriteString("\tsrc=\"" + mp3FileName + ".mp3\">\n")
	file.WriteString("\tYour browser does not support the <code>audio</code> element.\n")
	file.WriteString("</audio><br/>\n")

}

func writeHTMLImage(file *os.File, imageFileName string) {
	file.WriteString(" <img class=\"banner\" src=\"" + imageFileName + "\" class=\"resize\">")
}

func writeHTMLDiv(file *os.File) {
	file.WriteString("</div>\n")
}

func writeHTMLClose(file *os.File) {
	file.WriteString("</body>\n</html>\n")
}

func main() {
	var vocabWord string
	var vocabImages []string
	var targetImagePath string
	var reader, out *os.File
	var origImage, newImage image.Image

	//Create the HTML output file
	os.Remove("./html/test.html")
	htmlFile, err := os.Create("./html/test.html")
	defer htmlFile.Close()

	writeHTMLHeader(htmlFile)

	if len(os.Args) < 2 {
		fmt.Println("Please specify a file containing the vocab words")
		return
	}
	//get the name of the file with the vocab words
	vocabFile := os.Args[1]

	//Open the file with all the vocab words and flip out if the file doesn't exist
	fileHandle, err := os.Open(vocabFile)
	defer fileHandle.Close()
	check(err)
	scanner := bufio.NewScanner(fileHandle)

	//Scan through them 1 at a time...
	for scanner.Scan() {
		//vocabWord = scanner.Text()
		//vocabWord = strings.Join(makestring(scanner.Bytes()), "")
		vocabWord = scanner.Text()
		//if ! oxford.CheckIfExists(vocabWord) {
		if 1 == 2 {
			fmt.Println(vocabWord, " does not exist in Oxford Spanish Dictionary.  Check Spelling")
		} else {

			fmt.Println("processing: ", vocabWord)

			//TODO - Add Oxford Dictionary Check to make sure the word is not mis-spelled.

			//Grab a mp3 of the pronunciation and some jpegs from Bing image search

			//Get pronunciation from Forvo it it doesn't already exist
			if _, err := os.Stat("./html/" + vocabWord + ".mp3"); os.IsNotExist(err) {
				fmt.Println("Debug: before forvo search")
				forvosearch.GetPronuciation("es", vocabWord, "./html/")
				fmt.Println("Debug: after forvo search")

				fmt.Println("Debug:  mp3 File does not exist")
			} else {
				fmt.Println("Debug: mp3 file already exists")
			}

			writeHTMLSound(htmlFile, vocabWord)
			//Get images from Bing if they don't already exist
			numFiles, _ := filepath.Glob("./html/" + vocabWord + "*.jpg")
			fmt.Println("Debug: number of image files already existing:", len(numFiles))

			j := 0
			loopCount := 0
			if len(numFiles) < 9 {
				for j <= 8 {
					fmt.Println("Debug: before bing image search")
					vocabImages = bingImageSearch.ImageSearch(vocabWord)
					fmt.Println("Debug: after bing image search")

					loopCount++
					if loopCount > 20 {
						break
					}
					fmt.Println("Total:", j, "\tPass:", loopCount)
					fmt.Println(vocabImages[j])
					//fileExtension = path.Ext(vocabImages[j])
					//download the file from the URL to the ./html directory
					//TODO: Add a check for file extension .JPG or PNG before downloading
					targetImagePath = "./html/" + vocabWord + "_" + strconv.Itoa(j) + "_o.jpg"
					err = forvosearch.DownloadFile(targetImagePath, vocabImages[j])
					fmt.Println("Debug: Dowload Image URL error:", err)
					//Make all the images a consistent size of 150px by 120px
					//fmt.Println(targetImagePath)
					reader, err = os.Open(targetImagePath)
					check(err)

					fmt.Println("before Decode")
					origImage, _, err = image.Decode(reader)

					fmt.Println("Debug: Image Decode Error:", err)
					//fmt.Println("after Decode")
					reader.Close()
					if err == nil {
						//fmt.Println("before Resize")
						newImage = resize.Resize(150, 100, origImage, resize.Lanczos2)
						//fmt.Println("after Resize")
						out, err = os.Create("./html/" + vocabWord + "_" + strconv.Itoa(j) + ".jpg")
						check(err)
						jpeg.Encode(out, newImage, nil)
						out.Close()
						writeHTMLImage(htmlFile, vocabWord+"_"+strconv.Itoa(j)+".jpg")
						j++
					} else {
						//check(err)
					} //end for j <=8

				} //end-if number of files < 9
			} else {
				fmt.Println("Debug: Images already exist.  Skipping download")
				for k := 0; k <= 8; k++ {
					writeHTMLImage(htmlFile, vocabWord+"_"+strconv.Itoa(k)+".jpg")
				}

			}
			fmt.Println("DEBUG:  writing HTML")
			writeHTMLDiv(htmlFile)
		} // end-if word is in Dictionary

	} // end read from file

	writeHTMLClose(htmlFile)
	os.Remove("./html/*_o.jpg")

}
