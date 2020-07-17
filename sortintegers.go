package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"sync"
	"strconv"
	"sort"
)

const NumPartitions = 4

func main(){
	// get user integers
	userInputs := getUserInputs()

	// partition user inputs into 4 parts
	arrayPartitions := partition(userInputs)

	// sort partitions
	var wg sync.WaitGroup
	wg.Add(NumPartitions) // 1 goroutine per partition
	for i := 0 ; i < NumPartitions; i++ {
		// could also have used channels here 
		go sortPartition(&arrayPartitions[i], &wg) 
	}
	wg.Wait()

	// recursively merge partitions
	sortedArray := arrayPartitions[0]
	for _, sortedPartition := range arrayPartitions[1 : NumPartitions]{
		sortedArray = merge(sortedArray, sortedPartition)
	}

	// print 
	fmt.Println("Sorted Array: ", sortedArray)
}

func getUserInputs() []int {
	// slice of user integers with length 0 capacity 10
	UserIntegers := make([]int, 0, 10)

	// get user input
	fmt.Println("Input integers to be sorted e.g. '3 11 2 9 5'")
	inputReader := bufio.NewReader(os.Stdin)
	InputtedIntegers, _ := inputReader.ReadString('\n')

	// split inputted integers
	splitIntegers := strings.Fields(InputtedIntegers)

	// convert from string to integer
	for _, value := range splitIntegers {
		intValue, err := strconv.Atoi(value)
		if err == nil{
			UserIntegers = append(UserIntegers, intValue)
		}
	}
	return UserIntegers
}

func partition(userIntegers []int) [][]int {
	numIntegers := len(userIntegers)
	minPartitionSize := numIntegers / NumPartitions 
	numUnallocatedIntegers := numIntegers - minPartitionSize * NumPartitions
	arrayPartitions := [][]int {}
	partitionStart := 0

	for i := 0; i < NumPartitions; i++ {
		partitionEnd := partitionStart + minPartitionSize
		if numUnallocatedIntegers > 0 {
			partitionEnd++
			numUnallocatedIntegers--
		}
		arrayPartitions = append(arrayPartitions, userIntegers[partitionStart : partitionEnd])
		partitionStart = partitionEnd
	}
	return arrayPartitions
}

func sortPartition(partitionIntegers *[]int, wg *sync.WaitGroup) {
	fmt.Println("Unsorted partition: ", *partitionIntegers)
	sort.Ints(*partitionIntegers) 
	wg.Done()
}

func merge(leftArray, rightArray []int) []int{
	mergedSize := len(leftArray) + len(rightArray)
	mergedArray := make([]int, mergedSize, mergedSize)
	leftIndex, rightIndex := 0, 0

	for i := 0; i < mergedSize; i++ {
		// if left already all merged but right not
		if leftIndex >= len(leftArray) {
			// merge right
			mergedArray[i] = rightArray[rightIndex]
			rightIndex++
		// if right already all merged but left not
		} else if rightIndex >= len(rightArray) {
			// merge left
			mergedArray[i] = leftArray[leftIndex]
			leftIndex++
		// if left < right 
		} else if leftArray[leftIndex] < rightArray[rightIndex] {
			// merge left
			mergedArray[i] = leftArray[leftIndex]
			leftIndex++
		// if right > left 
		} else {
			// merge right
			mergedArray[i] = rightArray[rightIndex]
			rightIndex++
		}
	}
	return mergedArray
}