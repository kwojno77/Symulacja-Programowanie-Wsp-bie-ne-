package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type vertice struct {
	index          int
	myChanel       chan packet
	trapping       chan bool
	next           []*vertice
	passingPackets []int
	trap           bool
}

type packet struct {
	index           int
	visitedVertices []int
	health          int
}

func main() {
	//ARGUMENTY WEJŚCIOWE
	if len(os.Args) < 6 {
		return
	}
	justArgs := os.Args[1:]
	verticeNum, _ := strconv.Atoi(justArgs[0])
	if verticeNum < 3 {
		verticeNum = 3
	}
	shortcutNum, _ := strconv.Atoi(justArgs[1])
	if shortcutNum < 0 {
		shortcutNum = 0
	}
	backwardNum, _ := strconv.Atoi(justArgs[2])
	if backwardNum < 0 {
		backwardNum = 0
	}
	packetNum, _ := strconv.Atoi(justArgs[3])
	if packetNum < 1 {
		packetNum = 1
	}
	packetHealth, _ := strconv.Atoi(justArgs[4])
	if packetHealth < 1 {
		packetHealth = 1
	}

	//GENEROWANIE LICZB PSEUDOLOSOWYCH (ZIARNO)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	//Flaga kończąca symulację
	end := false
	//slice pakietów
	packets := make([]packet, packetNum)

	//TWORZENIE WIERZCHOŁKÓW I USTALENIE NASTĘPNIKÓW
	v := make([]vertice, verticeNum)
	v[verticeNum-1] = vertice{index: verticeNum - 1, myChanel: make(chan packet, 0), trapping: make(chan bool), passingPackets: make([]int, 0), trap: false}
	for i := verticeNum - 1; i > 0; i-- {
		v[i-1] = vertice{index: i - 1, myChanel: make(chan packet, 0), trapping: make(chan bool), passingPackets: make([]int, 0), trap: false}
		s := make([]*vertice, 0)
		v[i-1].next = append(s, &v[i])
	}

	//LOSOWANIE SKRÓTÓW I WPISANIE ICH JAKO NASTĘPNIKÓW
	for i := 0; i < shortcutNum; i++ {
		r := rand.New(rand.NewSource(r1.Int63()))
		k := r.Intn(verticeNum-1) + 1
		j := r.Intn(k)
		v[j].next = append(v[j].next, &v[k])
	}

	//LOSOWANIE KRAWĘDZI COFAJĄCYCH
	for i := 0; i < backwardNum; i++ {
		r := rand.New(rand.NewSource(r1.Int63()))
		j := r.Intn(verticeNum-3) + 2
		k := r.Intn(j-1) + 1
		v[j].next = append(v[j].next, &v[k])
	}

	//DRUKOWANIE GRAFU
	for i := 0; i < verticeNum-1; i++ {
		fmt.Print(v[i].index, "->")
		length := len(v[i].next)
		for j := 0; j < length; j++ {
			fmt.Print("(", v[i].next[j].index, ")")
		}
		fmt.Println("")
	}
	fmt.Println(v[verticeNum-1].index)
	fmt.Println("----------------------------------------")

	//WĄTEK ŹRÓDŁA	- INICJALIZACJA
	source := func(ve vertice) {
		length := len(ve.next)
		empty := make([]int, 0)
		for i := 0; i < packetNum; i++ {
			r := rand.New(rand.NewSource(r1.Int63()))

			//losowanie wierzchołka-odbiorcy
			luckyVertice := ve.next[r.Intn(length)]
			//pobranie kanału od wierzchołka-odbiorcy
			chanel := luckyVertice.myChanel

			//tworzenie pakietu
			p := packet{index: i}
			p.health = packetHealth

			//wpisanie loga do pakietu
			p.visitedVertices = append(empty, ve.index)
			//wpisanie loga do wierzchołka
			v[ve.index].passingPackets = append(v[ve.index].passingPackets, p.index)

			//wysłanie pakietu
			chanel <- p
			//tymczasowe uśpienie wierzhołka
			time.Sleep(time.Duration(r.Float64() * float64(time.Second)))
		}
	}

	//WĄTEK WIERZCHOŁKA - INICJALIZACJA
	verticeF := func(ve vertice, msg chan<- string, lost chan bool) {
		length := len(ve.next)
		veChanel := ve.myChanel
		for {
			r := rand.New(rand.NewSource(r1.Int63()))
			//odebranie pakietu
			select {
			case <-ve.trapping:
				ve.trap = true
				continue
			case p := <-veChanel:
				msg <- "pakiet " + strconv.Itoa(p.index) + " jest w wierzchołku " + strconv.Itoa(ve.index)

				//wpisanie loga do pakietu
				p.visitedVertices = append(p.visitedVertices, ve.index)
				//wpisanie loga do wierzchołka
				v[ve.index].passingPackets = append(v[ve.index].passingPackets, p.index)

				//liczba pozostałych transferów pakietu zmniejsza się
				p.health--
				//jeśli pakietowi nie został już żaden transfer, umiera
				if ve.trap == true {
					p.health = 0
				}

				if p.health < 1 {
					//jeśli wpadł w pułapkę, pułapka znika
					if ve.trap == true {
						ve.trap = false
						msg <- "pakiet " + strconv.Itoa(p.index) + " wpadł w pułapkę i umiera!"
					} else {
						msg <- "pakiet " + strconv.Itoa(p.index) + " umiera!"
					}
					packets[p.index] = p
					//wysyłana jest informacja o zagubionym pakiecie
					lost <- true
					continue
				}
				//tymczasowe uśpienie wierzhołka
				time.Sleep(time.Duration(r.Float64() * float64(time.Second)))
				sent := false

				for !sent {
					r = rand.New(rand.NewSource(r1.Int63()))
					//losowanie wierzchołka-odbiorcy
					luckyVertice := ve.next[r.Intn(length)]
					//pobranie kanału od wierzchołka-odbiorcy
					chanel := luckyVertice.myChanel
					//wysłanie pakietu
					select {
					case chanel <- p:
						sent = true
					case <-time.After(time.Second):
						continue
					}
				}
				//tymczasowe uśpienie wierzhołka
				time.Sleep(time.Duration(r.Float64() * float64(time.Second)))
			}
		}
	}

	//WĄTEK DRUKARZA - INICJALIZACJA
	printer := func(msg <-chan string) {
		for {
			fmt.Println(<-msg)
		}
	}

	//WĄTEK KŁUSOWNIKA - INICJALIZACJA
	poacher := func(vertices []vertice, msg chan<- string) {
		for end == false {
			//losowanie ilości czasu jaką kłusownik będzie spał
			r := rand.New(rand.NewSource(r1.Int63()))
			time.Sleep(time.Duration(r.Float64() * 3 * float64(time.Second)))
			//losowanie wierzchołka, w którym zostanie umieszczona pułapka
			unluckyVertice := r.Intn(len(vertices)-2) + 1
			msg <- "KŁUSOWNIK: zastawiam pułapkę w wierzchołku " + strconv.Itoa(unluckyVertice)
			//pułapką jest pakiet o indeksie -1
			vertices[unluckyVertice].trapping <- true
		}
	}

	printChan := make(chan string)
	lost := make(chan bool)

	go source(v[0])
	go printer(printChan)
	for i := 1; i < verticeNum-1; i++ {
		go verticeF(v[i], printChan, lost)
	}
	if len(v) > 2 {
		go poacher(v, printChan)
	}

	//WĄTEK UJŚCIA (WĄTEK GŁÓWNY)
	oChanel := v[verticeNum-1].myChanel
	for i := packetNum; i > 0; i-- {
		r := rand.New(rand.NewSource(r1.Int63()))
		select {
		case p := <-oChanel:
			printChan <- "pakiet " + strconv.Itoa(p.index) + " został odebrany"

			p.health--
			p.visitedVertices = append(p.visitedVertices, v[verticeNum-1].index)
			v[verticeNum-1].passingPackets = append(v[verticeNum-1].passingPackets, p.index)
			packets[p.index] = p
			time.Sleep(time.Duration(r.Float64() * float64(time.Second)))
		case <-lost:
			continue
		}
	}

	//GENEROWANIE RAPORTU
	fmt.Println("----------------------------------------")
	fmt.Println("RAPORT:")
	fmt.Println("   Pakiety:")
	for i := 0; i < packetNum; i++ {
		if packets[i].visitedVertices != nil {
			fmt.Println(packets[i].index, ": ", packets[i].visitedVertices)
		}
	}
	fmt.Println("   Wierzchołki:")
	for i := 0; i < verticeNum; i++ {
		fmt.Println(v[i].index, ": ", v[i].passingPackets)
	}
}
