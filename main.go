package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"syscall/js"
)

func adder(id string) func(val []js.Value) {
	sum := 0
	return func(val []js.Value) {
		sum++
		js.Global().Get("document").Call("getElementById", id).Set("innerText", fmt.Sprintf("%d", sum))
	}
}

func loadCallback() {
	add1 := adder("span1")
	js.Global().Set("add1", js.NewCallback(add1))

	// Part 4
	fn := func(val []js.Value) {
		go func() {
			js.Global().Get("document").Call("getElementById", "day6").Set("innerText", "")
			answ1, answ2 := parseInput(input, 10000)
			str := fmt.Sprintf("Day 6 solution was part1 %d, and part2 %d", answ1, answ2)
			js.Global().Get("document").Call("getElementById", "day6").Set("innerText", str)
		}()
	}
	js.Global().Set("day6", js.NewCallback(fn))

}

func main() {
	// Part 1
	js.Global().Get("document").Call("getElementById", "init").Set("innerText", "Done")
	// Part 2
	loadCallback()

	// Part3
	cb := js.NewEventCallback(js.PreventDefault, func(val js.Value) {
		n1 := js.Global().Get("document").Call("getElementById", "num1").Get("value").String()
		n2 := js.Global().Get("document").Call("getElementById", "num2").Get("value").String()

		x, _ := strconv.Atoi(n1)
		y, _ := strconv.Atoi(n2)
		//fmt.Println(x + y)
		js.Global().Get("document").Call("getElementById", "input1").Set("value", x+y)
	})
	js.Global().Get("document").Call("getElementById", "button1").Call("addEventListener", "click", cb)

	js.Global().Get("document").Call("getElementById", "part2").Set("hidden", "")

	js.Global().Get("document").Call("getElementById", "part3").Set("hidden", "")

	js.Global().Get("document").Call("getElementById", "part4").Set("hidden", "")

	c := make(chan struct{}, 0)
	<-c
}

// Part 4

type point struct {
	y, x, id, sum int
	isEdge        bool
}

func (a *point) distance(b point) int {
	return int(math.Abs(float64(a.x)-float64(b.x))) + int(math.Abs(float64(a.y)-float64(b.y)))
}

func (a *point) sumDist(mp points) (sum int) {
	for _, r := range mp {
		sum += a.distance(r)
	}
	return
}

type points []point

func parseInput(input string, part2Limit int) (part1, part2 int) {
	res := strings.Fields(input)
	mp := make(points, len(res)/2)
	for i, n := 0, 0; i < len(res); i = i + 2 {
		y, _ := strconv.Atoi(strings.TrimSpace(strings.Replace(res[i], ",", "", -1)))
		x, _ := strconv.Atoi(strings.TrimSpace(strings.Replace(res[i+1], ",", "", -1)))
		mp[n] = point{y: y, x: x, id: n}
		n++
	}
	maxY, maxX := getMaxYx(mp)

	// Left
	for y := 0; y < maxY; y++ {
		p := point{y: y, x: 0}
		sort.Slice(mp, func(i, j int) bool { return mp[i].distance(p) < mp[j].distance(p) })
		mp[0].isEdge = true
	}
	// Right
	for y := 0; y < maxY; y++ {
		p := point{y: y, x: maxX}
		sort.Slice(mp, func(i, j int) bool { return mp[i].distance(p) < mp[j].distance(p) })
		mp[0].isEdge = true
	}
	// Top
	for x := 0; x < maxX; x++ {
		p := point{y: 0, x: x}
		sort.Slice(mp, func(i, j int) bool { return mp[i].distance(p) < mp[j].distance(p) })
		mp[0].isEdge = true
	}
	// Bottom
	for x := 0; x < maxX; x++ {
		p := point{y: maxY, x: x}
		sort.Slice(mp, func(i, j int) bool { return mp[i].distance(p) < mp[j].distance(p) })
		mp[0].isEdge = true
	}

	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			p := point{y: y, x: x}
			if p.sumDist(mp) < part2Limit {
				part2++
			}
			sort.Slice(mp, func(i, j int) bool { return mp[i].distance(p) < mp[j].distance(p) })
			if !mp[0].isEdge && mp[0].distance(p) != mp[1].distance(p) {
				mp[0].sum++
			}
		}
	}

	sort.Slice(mp, func(i, j int) bool { return mp[i].sum > mp[j].sum })
	part1 = mp[0].sum
	return
}

func getMaxYx(input []point) (y, x int) {
	sort.Slice(input, func(i, j int) bool { return input[i].y > input[j].y })
	y = input[0].y
	sort.Slice(input, func(i, j int) bool { return input[i].x > input[j].x })
	x = input[0].x
	return
}

const input = `264, 340
308, 156
252, 127
65, 75
102, 291
47, 67
83, 44
313, 307
159, 48
84, 59
263, 248
188, 258
312, 240
59, 173
191, 130
155, 266
252, 119
108, 299
50, 84
172, 227
226, 159
262, 177
233, 137
140, 211
108, 175
278, 255
259, 209
233, 62
44, 341
58, 175
252, 74
232, 63
176, 119
209, 334
103, 112
155, 94
253, 255
169, 87
135, 342
55, 187
313, 338
210, 63
237, 321
171, 143
63, 238
79, 132
135, 113
310, 294
289, 184
56, 259`
