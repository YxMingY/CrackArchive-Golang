package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

/**
 *  The github.com/mholt/archiver/v3/rar.go  has been customized by me.
 */

var FileName string

var TheadCount = 4
var SelectedTestFunc func(passwd string, filename string) int
var NowPasswordLength = 1

func main() {
	fmt.Printf("Please input Archive File name:")
	fmt.Scanln(&FileName)
	_, err := os.Stat(FileName)
	if err != nil {
		fmt.Println("Failed to open file.:" + err.Error())
		exit()
		return
	}
	for {
		fmt.Printf("Select alphabet: Numbers[0], Alphas[1], Nums&Alphas[2]:")
		var betNum int
		fmt.Scanln(&betNum)
		if SelectBet(betNum) {
			break
		}
	}
	SelectedTestFunc = SelectTestFunc(FileName)
	if SelectedTestFunc == nil {
		fmt.Println("Only Support .zip or .rar file")
		fmt.Printf("Enter any key to exit.")
		fmt.Scanln()
	}

	startTime := time.Now().UnixMilli() //毫秒计时
	for NowPasswordLength <= 5 {
		fmt.Printf("\nSetting password length at %d\n", NowPasswordLength)
		//time.Sleep(time.Second)
		fmt.Printf("[Main]Has tested 00/10")
		count := 0
		var processChans []chan int
		var passSubmits []chan string
		var dones []chan struct{} // 遇到一个坑，如果用make初始化切片通道，各个通道会被初始化为nil
		var IPasswords = make([][]int8, TheadCount)
		var processes = make([]int, TheadCount)
		//testTime指的是次数
		testTimePerThread := int(math.Pow(float64(getSelectBetLength()), float64(NowPasswordLength)))/TheadCount + 1
		for i := 0; i < TheadCount; i++ { //为各个线程初始化通道，计算任务
			IPasswords[i] = make([]int8, NowPasswordLength)
			setIPasswordByNum(IPasswords[i], count)
			count += testTimePerThread
			processChans = append(processChans, make(chan int))
			passSubmits = append(passSubmits, make(chan string))
			dones = append(dones, make(chan struct{}))
		}
		for i := 0; i < TheadCount; i++ { //启动线程
			go DoCrack(i, IPasswords[i], testTimePerThread, processChans[i], passSubmits[i], dones[i])
		}
		for {
			for i := 0; i < TheadCount; i++ {
				select {
				case password := <-passSubmits[i]:
					fmt.Printf("\n[Main]Found password %s\n", password)
					fmt.Printf("[Main]Used %.03f seconds.\n", float64(time.Now().UnixMilli()-startTime)/1000)
					KillThreads(dones)
					exit()
					return
				case processes[i] = <-processChans[i]:
				default:
				}
			}

			//进度计算与输出
			sumProcess := 0
			echoProg := true
			for i := 0; i < TheadCount; i++ {
				if processes[i] <= 0 {
					echoProg = false
					break
				}
				sumProcess += processes[i]
			}
			if echoProg {
				fmt.Printf("\r[Main]Has tested %02d/10", sumProcess/TheadCount)
			}

			//判断所有线程结束后跳出循环
			allComplete := true
			for i := 0; i < TheadCount; i++ {
				if processes[i] < 10 {
					allComplete = false
				}
			}
			if allComplete {
				KillThreads(dones)
				break
			}
		}
		NowPasswordLength++
	}
	fmt.Printf("\n[Main]Used %.03f seconds.\n", float64(time.Now().UnixMilli()-startTime)/1000)
}

// DoCrack
// @parameter
//
//	IPassword
//	testTime: Do how many times of test
//	process: used to report progress
//	passSubmit: When found the password, send
//	done: Used by main to shut down the routine
func DoCrack(TId int, iPassword []int8, testTime int, process chan int, passSubmit chan string, done chan struct{}) {
	prog := 0
	i := 0
	for {
		select {
		case <-done:
			//fmt.Println("\nThread has been killed")
			return
		default:
			password, end := NextPasswd(iPassword)
			i++
			if end {
				i = testTime //为了解决计算结果整数导致i可能达不到testTime的问题
			}
			if int(float32(i+1)/float32(testTime)*10) > prog {
				prog = int(float32(i+1) / float32(testTime) * 10)
				//fmt.Println(TId, i, prog, password)
				go func() {
					process <- prog //进度传回主线程,不能让自己被阻塞
				}()
			}
			if SelectedTestFunc(password, FileName) == 1 {
				fmt.Printf("\n[Thread %d]Found Password:%s\n", TId, password)
				passSubmit <- password //找到密码发回主线程并等待叫停
				<-done
				return
			}
			if i == testTime {
				<-done //任务做完就等主线程叫停
				return
			}
		}
	}
}

func KillThreads(dones []chan struct{}) {
	for i := 0; i < TheadCount; i++ {
		dones[i] <- struct{}{}
	}
}

func exit() {
	fmt.Printf("Enter any key to exit.")
	fmt.Scanln()
}
