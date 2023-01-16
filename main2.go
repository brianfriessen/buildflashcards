package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	bingImageSearch "example.com/bingImageSearch"
	forvosearch "example.com/forvosearch"
	"github.com/nfnt/resize"
)

const NUM_IMAGES = 12

var wg sync.WaitGroup
var decodeWG sync.WaitGroup

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

func DownloadFile(filepath string, url string) error {
	defer wg.Done()
	fmt.Println(url)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func decodeImage(htmlFile *os.File, origImagePath string, j int, vocabWord string) {
	defer decodeWG.Done()
	fmt.Println("Debug:", origImagePath)
	fmt.Println("before Resize")
	reader, err := os.Open(origImagePath)
	origImage, _, err := image.Decode(reader)
	if err == nil {
		fmt.Println("Debug: Open image error is", err)
		newImage := resize.Resize(150, 100, origImage, resize.Lanczos2)
		fmt.Println("after Resize")
		fmt.Println("Debug before Decode")
		out, err := os.Create("./html/" + vocabWord + "_" + strconv.Itoa(j) + ".jpg")
		fmt.Println("Debug: Image Decode Error:", err)
		check(err)
		jpeg.Encode(out, newImage, nil)
		out.Close()
		//writeHTMLImage(htmlFile, vocabWord+"_"+strconv.Itoa(j)+".jpg")

	}

}

func downloadImages(targetWord string) {
	//defer wg.Done()
	vocabImages := bingImageSearch.ImageSearch(targetWord)
	startIndex := 0
	stopIndex := startIndex + NUM_IMAGES
	results, _ := filepath.Glob("./html/" + targetWord + "_*o.jpg")
	fmt.Println("Number of files is: ", len(results), "\n\n\n")
	for len(results) < NUM_IMAGES {
		for j := startIndex; j < stopIndex; j++ {
			wg.Add(1)
			fmt.Println(vocabImages[j])
			targetImagePath := "./html/" + targetWord + "_" + strconv.Itoa(j) + "_o.jpg"
			go DownloadFile(targetImagePath, vocabImages[j])

		}
		wg.Wait()
		startIndex = stopIndex + 1
		stopIndex = startIndex + NUM_IMAGES
		results, _ := filepath.Glob("./html/" + targetWord + "_*o.jpg")
		if len(results) >= NUM_IMAGES {
			break
		}
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

//TODO: DEBUG Negative Wait Group when files are already there and we are skipping

func main() {
	var vocabWord string

	//Basic error checking
	if len(os.Args) < 3 {
		fmt.Println("Please specify a file containing the vocab words and spanish or german")
		return
	}

	//get the name of the file with the vocab words
	vocabFile := os.Args[1]
	langType := os.Args[2]

	//Create the HTML output file
	os.Remove("./html/" + vocabFile + ".html")
	htmlFile, err := os.Create("./html/" + vocabFile + ".html")
	check(err)
	defer htmlFile.Close()

	//Add the HTML Header Info---doesn't change
	writeHTMLHeader(htmlFile)

	//Open the file specified on the command line with all the vocab words and flip out if the file doesn't exist
	fileHandle, err := os.Open(vocabFile)
	defer fileHandle.Close()
	check(err)
	scanner := bufio.NewScanner(fileHandle)

	//Scan through them 1 at a time...
	for scanner.Scan() {
		/* Not sure why this was needed before but now bufio scanner works with UTF-8 accent and tilde chars.
		vocabWord = scanner.Text()
		vocabWord = strings.Join(makestring(scanner.Bytes()), "")
		*/
		vocabWord = scanner.Text()

		//if ! oxford.CheckIfExists(vocabWord) { TODO: add some delay in oxford check do we don't bomb out the API
		if 1 == 2 {
			fmt.Println(vocabWord, " does not exist in Oxford Spanish Dictionary.  Check Spelling")
		} else {

			fmt.Println("processing: ", vocabWord)

			//TODO - Add Oxford Dictionary Check to make sure the word is not mis-spelled.
			duration := time.Second
			time.Sleep(duration)
			//Get pronunciation from Forvo it it doesn't already exist
			if _, err := os.Stat("./html/" + vocabWord + ".mp3"); os.IsNotExist(err) {
				fmt.Println("Debug: before forvo search")
				if langType == "spanish" {
					forvosearch.GetPronuciation("es", vocabWord, "./html/")

				} else {
					forvosearch.GetPronuciation("de", vocabWord, "./html/")
				}
				fmt.Println("Debug: after forvo search")
				fmt.Println("Debug:  mp3 File does not exist")
			} else {
				fmt.Println("Debug: mp3 file already exists")
			}

			writeHTMLSound(htmlFile, vocabWord)

			//Get images from Bing if they don't already exist
			numFiles, _ := filepath.Glob("./html/" + vocabWord + "*.jpg")
			fmt.Println("Debug: number of image files already existing:", len(numFiles))

			if len(numFiles) < 9 {

				downloadImages(vocabWord) // func should return NUM_IMAGES JPEG/PNG files in the html dir

				decodeWG.Add(NUM_IMAGES) // Reformat all the images concurrently
				allImageFiles, _ := filepath.Glob("./html/" + vocabWord + "*.jpg")
				for j := range allImageFiles {
					go decodeImage(htmlFile, allImageFiles[j], j, vocabWord)
				}
				decodeWG.Wait()

			}

			//should now have NUM_IMAGES scaled image files in html dir
			for j := 0; j < NUM_IMAGES; j++ {
				writeHTMLImage(htmlFile, vocabWord+"_"+strconv.Itoa(j)+".jpg")
			}
		} // end-if word is in Dictionary

		fmt.Println("DEBUG:  writing HTML")
		writeHTMLDiv(htmlFile)
	} // end read from file

	writeHTMLClose(htmlFile)

	//clean-up files:
	cwd, _ := syscall.Getwd()
	fmt.Println(cwd)
	removeName := cwd + "./html/*_o.jpg"
	fmt.Println(removeName)
	results, _ := filepath.Glob(removeName)
	fmt.Println(len(results))
	for i := range results {
		fmt.Println("removing: ", results[i])
		os.Remove(results[i])
	}

}
