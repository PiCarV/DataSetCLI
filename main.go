package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

func main() {
    fmt.Println(`
     _______              __                 ______              __       ______   __        ______ 
    |       \            |  \               /      \            |  \     /      \ |  \      |      \
    | $$$$$$$\  ______  _| $$_     ______  |  $$$$$$\  ______  _| $$_   |  $$$$$$\| $$       \$$$$$$
    | $$  | $$ |      \|   $$ \   |      \ | $$___\$$ /      \|   $$ \  | $$   \$$| $$        | $$  
    | $$  | $$  \$$$$$$\\$$$$$$    \$$$$$$\ \$$    \ |  $$$$$$\\$$$$$$  | $$      | $$        | $$  
    | $$  | $$ /      $$ | $$ __  /      $$ _\$$$$$$\| $$    $$ | $$ __ | $$   __ | $$        | $$  
    | $$__/ $$|  $$$$$$$ | $$|  \|  $$$$$$$|  \__| $$| $$$$$$$$ | $$|  \| $$__/  \| $$_____  _| $$_ 
    | $$    $$ \$$    $$  \$$  $$ \$$    $$ \$$    $$ \$$     \  \$$  $$ \$$    $$| $$     \|   $$ \
     \$$$$$$$   \$$$$$$$   \$$$$   \$$$$$$$  \$$$$$$   \$$$$$$$   \$$$$   \$$$$$$  \$$$$$$$$ \$$$$$$
                                                                                                    
    `)
    var input = flag.String("input", "" ,"Input path to dataset")
    var output = flag.String("output", "" ,"Output path to dataset")
    var shuffle = flag.Bool("shuffle", false, "Shuffle dataset")
    var split = flag.Float64("split", 1.0, "Split dataset if 1 then all data is used for training, if 0.8 then 80% of the data is used for training and 20% is used for testing")
    flag.Parse()

    // check that input and output are not empty
    if *input == "" || *output == "" {
        fmt.Println("Input and output paths are required")
        return
    }

    // check that split is between 0 and 1
    if *split < 0 || *split > 1 {
        fmt.Println("Split must be between 0 and 1")
        return
    }

    // check that shuffle is true or false
    if *shuffle != true && *shuffle != false {
        fmt.Println("Shuffle must be true or false")
        return
    }


    // check that the input path is a directory that contains at least one jpg file and one csv file called steeringData.csv
    fmt.Println("Checking input path...")

    if !fileExists(*input + "/steeringData.csv") {
        fmt.Println("Input path does not exist or is missing steeringData.csv")
        os.Exit(1)
        return
    }

    fmt.Println("Input path is valid")

    // check that the output path is a directory
    fmt.Println("Checking output path...")
    if !directoryExits(*output) {
        fmt.Println("Output path does not exist")
        os.Exit(2)
        return
    }

    // check that the output path is empty
    fmt.Println("Checking output path is empty...")
    empty, err := IsEmpty(*output)
    if err != nil {
        fmt.Println("Error checking if output path is empty")
        os.Exit(3)
        return
    }
    if !empty {
        fmt.Println("Output path is not empty")
        os.Exit(4)
        return
    }


    // loop through all files ending in .jpg in the input directory and store them in a variable
    fmt.Println("Reading input directory...")
    var items, _ = ioutil.ReadDir(*input)
    var jpgs []string 
    for _, item := range items {
        if item.Name()[len(item.Name())-4:] == ".jpg" {
            jpgs = append(jpgs, item.Name())
        }
         if item.Name()[len(item.Name())-4:] == ".csv" {
            println("csv file found")
         }
    }

    // sort the jpgs by there names numerically
    jpgs = sortJpgs(jpgs)


    // create a new directory in the output path called data
    fmt.Println("Creating data directory...")
    err = os.Mkdir(*output + "/data", 0755)
    if err != nil {
        fmt.Println("Error creating data directory")
        os.Exit(5)
        return
    }

    // create a new directory in the output path called data/train
    fmt.Println("Creating train directory...")
    err = os.Mkdir(*output + "/data/train", 0755)
    if err != nil {
        fmt.Println("Error creating train directory")
        os.Exit(6)
        return
    }

    // create a new directory in the output path called data/test
    fmt.Println("Creating test directory...")
    err = os.Mkdir(*output + "/data/test", 0755)
    if err != nil {
        fmt.Println("Error creating test directory")
        os.Exit(7)
        return
    }

    // create and open 2 new csv files in the output path called data/train/steeringData.csv and data/test/steeringData.csv
    fmt.Println("Creating train csv...")
    trainCsv, err := os.Create(*output + "/data/train/steeringData.csv")
    if err != nil {
        fmt.Println("Error creating train csv")
        os.Exit(8)
        return
    }

    fmt.Println("Creating test csv...")
    testCsv, err := os.Create(*output + "/data/test/steeringData.csv")
    if err != nil {
        fmt.Println("Error creating test csv")
        os.Exit(9)
        return
    }

    // open the csv file in the input path called steeringData.csv
    fmt.Println("Opening input csv...")
    csvFile, err := os.Open(*input + "/steeringData.csv")
    if err != nil {
        fmt.Println("Error opening input csv")
        os.Exit(10)
        return
    }

    // read the csv file
    fmt.Println("Reading input csv...")
    inputScanner := bufio.NewScanner(csvFile)

    // count the number of lines in the csv file
    fmt.Println("Counting lines in input csv...")
    var lineCount int
    for inputScanner.Scan() {
        if inputScanner.Text() != "" {
            lineCount++
        }
    }
    fmt.Println("Number of lines in input csv:", lineCount)

    // check if the number of lines in the csv file is the same as the number of jpgs in the input directory
    if lineCount != len(jpgs) {
        fmt.Println("Number of lines in csv does not match number of jpgs in input directory")
        os.Exit(12)
        return
    }

    // // take the number of lines in the csv file and multiply it by the split value to get the number of lines to put in the train csv
    var trainLineCount = int(float64(lineCount) * *split)

    // // take the number of lines in the csv file and subtract the number of lines to put in the train csv to get the number of lines to put in the test csv
    var testLineCount = lineCount - trainLineCount

    // check that the train and test line counts are equal to the total line count
    if trainLineCount + testLineCount != lineCount {
        // this should never happen
        fmt.Println("Error calculating train and test line counts")
        os.Exit(13)
        return
    }

    // copy the first trainLineCount lines from the input csv to the train csv
    fmt.Println("Copying train lines to train csv...")
    csvFile.Seek(0, 0)
    inputScanner = bufio.NewScanner(csvFile)

    bufioWriter := bufio.NewWriter(trainCsv)
    var outputCount = 0
    for inputScanner.Scan() {
        if outputCount == trainLineCount {
            fmt.Println("Train line count reached at" , outputCount)
            break
        }
        _, err = bufioWriter.WriteString(inputScanner.Text())
        if err != nil {
            fmt.Println("Error writing to train csv")
            os.Exit(14)
            return
        }
        _, err = bufioWriter.WriteString("\n")
        if err != nil {
            fmt.Println("Error writing to train csv")
            os.Exit(15)
            return
        }
        outputCount++
    }


    // close the train csv
    bufioWriter.Flush()
    trainCsv.Close()

    fmt.Println("Finished copying train lines to train csv")

    // create a new bufio writer for the test csv
    bufioWriter = bufio.NewWriter(testCsv)

    // continue reading the input csv until the end of the testlinecount
    fmt.Println("Copying test lines to test csv...")
    outputCount = 0
    for inputScanner.Scan() {
        if outputCount == testLineCount {
            fmt.Println("Test line count reached at" , outputCount)
            break
        }
        _, err = bufioWriter.WriteString(inputScanner.Text())
        if err != nil {
            fmt.Println("Error writing to test csv")
            os.Exit(16)
            return
        }
        _, err = bufioWriter.WriteString("\n")
        if err != nil {
            fmt.Println("Error writing to test csv")
            os.Exit(17)
            return
        }
        outputCount++
    }


    // close the test csv
    bufioWriter.Flush()
    testCsv.Close()


    fmt.Println("Finished copying test lines to test csv")


// start copying the jpgs to the train and test directories
    fmt.Println("Copying jpgs to train and test directories...")

    for i, jpg := range jpgs {
        if i < trainLineCount {
            _, err = copy(*input + "/" + jpg, *output + "/data/train/" + jpg)
            if err != nil {
                fmt.Println("Error copying jpgs")
                os.Exit(18)
                return
            }
    }
        if i >= trainLineCount {
            _, err = copy(*input + "/" + jpg, *output + "/data/test/" + jpg)
        }
        if err != nil {
            fmt.Println("Error copying jpgs")
            os.Exit(19)
            return
        }
    }




    






    




}


func sortJpgs(jpgs []string) []string {
    sort.Slice(jpgs, func(i, j int) bool {
        return jpgs[i] < jpgs[j]
    })
    return jpgs
}


func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
       return false
    }
    return !info.IsDir()
 }

 func directoryExits(dir string) bool {
    info, err := os.Stat(dir)
    if os.IsNotExist(err) {
       return false
    }
    return info.IsDir()
 }

 func IsEmpty(name string) (bool, error) {
    f, err := os.Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1) // Or f.Readdir(1)
    if err == io.EOF {
        return true, nil
    }
    return false, err // Either not empty or error, suits both cases
}

func copy(src, dst string) (int64, error) {
    sourceFileStat, err := os.Stat(src)
    if err != nil {
            return 0, err
    }

    if !sourceFileStat.Mode().IsRegular() {
            return 0, fmt.Errorf("%s is not a regular file", src)
    }

    source, err := os.Open(src)
    if err != nil {
            return 0, err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
            return 0, err
    }
    defer destination.Close()
    nBytes, err := io.Copy(destination, source)
    return nBytes, err
}