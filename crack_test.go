package main

import (
	"fmt"
	"math"
	"testing"
)

func BenchmarkCrack(b *testing.B) {
	for ii := 0; ii < b.N; ii++ {
		FileName = "m.zip"
		SelectBet(0)
		SelectedTestFunc = SelectTestFunc(FileName)
		for NowPasswordLength <= 5 {
			fmt.Printf("\nSetting password length at %d\n", NowPasswordLength)
			//time.Sleep(time.Second)
			//fmt.Printf("[Main]Has tested 00/10")
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
						//fmt.Printf("[Main]Used %.03f seconds.\n", float64(time.Now().UnixMilli()-startTime)/1000)
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
					//fmt.Printf("\r[Main]Has tested %02d/10", sumProcess/TheadCount)
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
	}
}
