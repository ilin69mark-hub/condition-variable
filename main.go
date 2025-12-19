package main

import (
	"fmt"
	"sync"
)

const step int64 = 1
const targetIteration int = 1000 // Целевое количество итераций

func main() {
	var counter int64 = 0
	var mutex sync.Mutex
	var cond = sync.NewCond(&mutex)
	
	// Целевое значение счетчика
	targetValue := int64(targetIteration * int(step))
	
	// Функция инкремента
	increment := func(wg *sync.WaitGroup, iterationNum int) {
		defer wg.Done()
		
		mutex.Lock()
		counter += step
		
		// Проверяем, достиг ли счетчик целевого значения
		if counter >= targetValue {
			cond.Broadcast() // Уведомляем все ожидаемые горутины
		}
		mutex.Unlock()
	}
	
	// Запускаем горутины
	var wg sync.WaitGroup
	for i := 1; i <= targetIteration; i++ {
		wg.Add(1)
		go increment(&wg, i)
	}
	
	// Ждем, пока счетчик достигнет целевого значения
	mutex.Lock()
	for counter < targetValue {
		cond.Wait() // Ожидаем, пока не будет достигнуто целевое значение
	}
	mutex.Unlock()
	
	// Ждем завершения всех горутин
	wg.Wait()
	
	fmt.Printf("Счетчик достиг значения: %d\n", counter)
	fmt.Printf("Целевое значение: %d\n", targetValue)
}