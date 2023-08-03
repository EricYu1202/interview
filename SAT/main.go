package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	_stop = iota
	_beef
	_pork
	_chicken
)

const _workerCount = 5

var _wg = &sync.WaitGroup{}

func main() {
	//所有肉品
	production := []int{_beef, _beef, _beef, _beef, _beef, _beef, _beef, _beef, _beef, _beef,
		_pork, _pork, _pork, _pork, _pork, _pork, _pork,
		_chicken, _chicken, _chicken, _chicken, _chicken}

	//打亂肉品順序
	randomSeed := rand.New(rand.NewSource(time.Now().Unix()))
	randomSeed.Shuffle(len(production), func(i, j int) {
		production[i], production[j] = production[j], production[i]
	})

	//放入停止信號
	production = append(production, []int{_stop, _stop, _stop, _stop, _stop}...)

	productionLine := make(chan int, len(production))
	_wg.Add(1)
	go func() { //放到流水線上
		defer _wg.Done()
		for _, value := range production {
			productionLine <- value
		}
		close(productionLine)
	}()
	_wg.Add(1)
	go work(_workerCount, productionLine)

	_wg.Wait()

}

// 工廠運作
func work(count int, productionLine <-chan int) {
	defer _wg.Done()
	for i := 0; i < count; i++ { //總共5個工人
		_wg.Add(1)
		go func(worker string) {
			defer _wg.Done()
			for {
				meat := <-productionLine //一有肉品放到channel上就馬上處理
				if meat != _stop {
					handle(worker, meat)
				} else {
					break
				}
			}

		}(string(byte(i + 65)))
	}
}

// 工人處理肉品
func handle(worker string, meat int) {

	switch meat {
	case _beef:
		log.Printf("%s 在 %s 取得牛肉\n", worker, time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(5 * time.Second)
		log.Printf("%s 在 %s 處理完牛肉\n", worker, time.Now().Format("2006-01-02 15:04:05"))
	case _pork:
		log.Printf("%s 在 %s 取得豬肉\n", worker, time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(2 * time.Second)
		log.Printf("%s 在 %s 處理完豬肉\n", worker, time.Now().Format("2006-01-02 15:04:05"))
	case _chicken:
		log.Printf("%s 在 %s 取得雞肉\n", worker, time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(3 * time.Second)
		log.Printf("%s 在 %s 處理完雞肉\n", worker, time.Now().Format("2006-01-02 15:04:05"))
	}
}
