package wb_internship_l1_theory

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"
	"unsafe"
)

func question1() {
	/*
		1. Какой самый эффективный способ конкатенации строк?
	*/
	// Go < 1.10
	buffer := bytes.Buffer{}
	buffer.WriteString("123")
	// Go > 1.10
	builder := new(strings.Builder)
	builder.WriteString("123")
}

/*
	2. Что такое интерфейсы, как они применяются в Go?
*/
type Player interface {
	Play()
}
type Human struct {
	// какие нибудь поля
}

// Куча методов для человека
func (h *Human) Play() {
	fmt.Println("Human is playing xbox")
}

type Dog struct {
	// какие нибудь другие поля, разные методы, но тоже есть play
}

func (d *Dog) Play() {
	fmt.Println("Dog is playing with ball")
}

func interrupt(p Player) {
	fmt.Println("Play time is over!")
}

func question2() {
	// Интерфейсы - типы, отражающие реализуемые структурой методы
	// Чтобы структура соответствовала интерфейсу - она должна реализовывать все методы интерфейса
	h := &Human{}
	d := &Dog{}
	interrupt(h)
	interrupt(d)
}

/*
Чем отличаются RWMutex от Mutex?
*/

func question3() {
	mu := sync.Mutex{} // Блокирует доступ к переменной для горутин для того,
	// Чтобы конкурентное чтение / запись не вызвало проблем
	someNumber := 10
	wg := sync.WaitGroup{}
	wg.Add(1)
	mu.Lock()
	go func() {
		fmt.Println("Блокируем mu на 5 секунд. Следующая горутина должна читать значение!")
		time.Sleep(5 * time.Second)
		someNumber = 12
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		mu.Lock()
		fmt.Println(someNumber)
		mu.Unlock()
	}()
	wg.Wait()
	//rwMu := sync.RWMutex{} // Блокирует доступ к переменной для горутин. Lock, Unlock - управляет записью
	//// RLock, RUnlock - чтение. Заблокированная для записи переменная все еще доступна для чтения другими
	//// горутинами
	rwMu := sync.RWMutex{}
	builder := new(strings.Builder)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			rwMu.Lock()
			time.Sleep(1 * time.Second)
			builder.WriteString("ы")
			rwMu.Unlock()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			time.Sleep(500 * time.Millisecond)
			rwMu.RLock()
			fmt.Println(builder) // Можем читать ее.
			rwMu.RUnlock()
		}
	}()
	wg.Wait()
}

/*
	4. Чем отличаются буферизированные и не буферизированные каналы?
*/
// Буфферезированные каналы не блокируют выполнение горутины, пока их буффер не заполнен.
//Небуфферезированные каналы блокируют канал до отправки/получения сообщения

/*
	5.Какой размер у структуры struct{}{}?
*/

type LolType struct{}

func (l *LolType) hello() {
	fmt.Println("hello world")
}
func question5() {
	lol := struct{}{}
	var lol2 LolType
	fmt.Println(unsafe.Sizeof(lol))  // 0!
	fmt.Println(unsafe.Sizeof(lol2)) // 0!
	lol2.hello()
}

/*
6. Есть ли в Go перегрузка методов или операторов?
Ответ: Нету!
*/

/*
7.В какой последовательности будут выведены элементы map[int]int?

Пример:
m[0]=1
m[1]=124
m[2]=281
Ответ: Вывод не может быть гарантирован
*/
func question7() {
	newMap := make(map[int]int)
	newMap[0] = 1
	newMap[1] = 124
	newMap[2] = 201
	for k, v := range newMap {
		fmt.Println(k, v)
	}
}

/*
8. В чем разница make и new?
new выделяет память и возвращает указатель. new работает для любого типа
make инициализирует тип, возвращая значение. Работает с внутренними типами slice, map, chan
*/

/*
9. Сколько существует способов задать переменную типа slice или map?
*/
func question9() {
	// Slice
	var slice1 = []string{"123", "234", "456"}
	slice2 := []string{"123", "234", "456"}
	slice3 := make([]string, 3)
	arr := [3]string{"1", "2", "3"}
	slice4 := arr[:]
	slice5 := new([]string)
	*slice5 = append(*slice5, "1", "2", "3")
	fmt.Println(slice1, slice2, slice3, slice4, *slice5)
	// Map
	var map1 = map[int]int{1: 2, 3: 4, 5: 6}
	map2 := map[int]int{1: 1}
	map3 := make(map[int]int)
	map3[2] = 22
	fmt.Println(map1, map2, map3)
}

/*
10. Что выведет данная программа и почему?
func update(p *int) {
 // здесь копия p, но тк это указатель, он все еще указывает на a из main
  b := 2 // создание новой переменной
  p = &b // теперь МЕСТНЫЙ p указывает на b. Этот p исчезнет с выходом из области видимости....
}

func main() {
  var (
     a = 1 // новая переменная
     p = &a // указатель на a
  )
  fmt.Println(*p) // значение a по указателю
  update(p) // передача указателя
  fmt.Println(*p) // значение исходного а по указателю не поменялось, поэтому здесь будет 1
}
*/

/*
11.Что выведет данная программа и почему?

func main() {
  wg := sync.WaitGroup{}
  for i := 0; i < 5; i++ {
     wg.Add(1)
     go func(wg sync.WaitGroup, i int) { // передаем значение, поэтому ->
        fmt.Println(i)
        wg.Done() -> значение, измененное здесь не повлияет на wg, объявленный выше горутины
		// Здесь wg закрывается, но во внешней области видимости он все еще будет равен 5
     }(wg, i)
  }
	// получаем фатал здесь, потому что не дождемся. Надо было юзать вг из внешней области (не передавать параметром)
	// или передавать поинтер
  wg.Wait()
  fmt.Println("exit")
}

*/

/*
12. Что выведет данная программа и почему?
func main() {
  n := 0
  if true {
     n := 1 // новая переменная, которая актуальна только в этой области видимости
	// Плохая практика, IDE наверняка будет ругаться, но теперь это новая переменая, которая доживет до }
     n++
  }
  fmt.Println(n) // 0
	// Если хотим 1, надо поменять := на =
}
*/

/*
13.Что выведет данная программа и почему?
func someAction(v []int8, b int8) {
	v[0] = 100
	v = append(v, b) // но здесь меняется уже сам массив, а значение сохраняется только
	// в этой функции. 0 индекс изменится в исходом массиве, но 6ому там взяться неоткуда
}

func main() {
	var a = []int8{1, 2, 3, 4, 5}
	someAction(a, 6) // слайс всего лишь отображает значения лежащего под ней массива,
	// поэтому если слайс ссылается на один и тот же массив, он изменится при изменении
	// в другой функции
	fmt.Println(a)
}
*/

/*
14.Что выведет данная программа и почему?
func main() {
  slice := []string{"a", "a"}

  func(slice []string) {
     slice = append(slice, "a") // теперь под slice другой массив. Это другое значение!
     slice[0] = "b"
     slice[1] = "b"
     fmt.Print(slice) // b, b, a
  }(slice) // тоже самое, если бы функция была отдельной. Передаем как параметр значение
  fmt.Print(slice) // a, a
}

*/
