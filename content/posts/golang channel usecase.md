---
title: "Một vài cách sử dụng channel trong golang"
date: 2023-03-03
draft: false
authors: ["ngoctd"]
tags: ["golang", "channel", "concurrency"]
---

Lập trình bất đồng bộ (asynchronous) và đồng thời (concurrency) trở  lên dễ dàng với channel trong Golang. Synchronization với channel có phạm vi sử dụng lớn và nhiều biến thể hơn so với các mô hình khác như actoc model hay async/await pattern.

Trong bài này mình sẽ tổng hợp lại các trường hợp sử dụng channel trong golang, mục đích chính là sử dụng nhiều nhất có thể (chứ không phải là trường hợp tốt nhất) có thể các kĩ thuật khác về synchronization sẽ tốt hơn.

## Sử dụng channel như Futures/Promises

Future và promises được sử dụng trong nhiều ngôn ngữ khác nhau. Chúng thường được sử dụng dưới dạng requests và responses.

**Lười làm thì thuê - Return receive-only channel**

Trong ví dụ dưới longTimeRequest xử lý async bằng cách trả về receive-only channel. Mỗi hàm longTimeRequest thực hiện hết 1s (sleep). Bằng cách đưa phần sleep (minh họa cho workload nào đó) chạy dưới goroutine và trả về ngay receive-only channel thì giúp x2 thời gian thực hiện.

```go
func main() {
	now := time.Now()

    // do async
	r1, r2 := longTimeRequest(), longTimeRequest()

	getResponse(r1)
	getResponse(r2)

	log.Println("since: ", time.Since(now))
}

func longTimeRequest() <-chan string {
	ch := make(chan string)
	now := time.Now()

	// async here
	go func() {
		time.Sleep(time.Second)
		ch <- fmt.Sprintf("long time request response after: %v", time.Since(now))
        close(ch)
	}()

	return ch
}

func getResponse(ch <-chan string) {
	log.Println(<-ch)
}
```

**Lười làm thì thuê - Pass send-only channels**

Thay vì trả về receive-only channel, ta có thể sử dụng send-only channels như là một phương pháp thay thế. Tuy nhiên mình thấy cách receive-only channel sử dụng clean hơn. 

```go
func main() {
	now := time.Now()

    // sử dụng buffered channel để tránh block khi gửi
	ch := make(chan string, 2)
	go longTimeRequest(ch)
	go longTimeRequest(ch)

	getResponse(ch)
	getResponse(ch)

	log.Println("since: ", time.Since(now))
}

func longTimeRequest(ch chan<- string) {
	now := time.Now()
	time.Sleep(time.Second)
	ch <- fmt.Sprintf("long time request response after: %v", time.Since(now))
}

func getResponse(ch <-chan string) {
	log.Println(<-ch)
}
```

**Ai nhanh thì thắng - The first response win**

Một ngày đẹp trời, cà rốt của bạn đăng story: tối nay tớ rảnh? Bạn suy nghĩ hết nửa tiếng đồng hồ rồi reply để hẹn địa điểm, thì nhận được câu trả lời là tớ set kèo đi với người khác rồi. Thực ra là dù có reply sớm thì cà rốt cũng chả đi với bạn đâu :v tuy nhiên bạn đến muộn thì người ta có cái cớ rất chi là hợp lí để từ chối. Không lan man nữa, vào việc thôi.

```go
func main() {
	now := time.Now()

	queue := 10
	ok := make(chan int, queue)
	for i := 0; i < queue; i++ {
		go reply(ok, i)
	}

	// ai nhanh thì thắng
	firstRepsonse := <-ok
	log.Println("since: ", time.Since(now), "pos", firstRepsonse)
}

func init() {
	rand.Seed(time.Now().Unix())
}

func reply(ch chan<- int, pos int) {
	rd := rand.Intn(3) + 1
	time.Sleep(time.Duration(rd) * time.Second)

	ch <- pos
}
```

Chú ý là phải dùng buffered channel để tránh ăn send block của channel nhé (goroutine leak problem). Đoạn code trên còn vấn đề là chẳng hạn cà rốt chỉ đăng tin cho người nào đó xem thôi, người đó reply rồi thì cà rốt thêm một tin mới là đã tìm được người cần tìm rồi. Thế là bạn cũng không cần reply lại nữa. Vậy xử lý thế nào nhỉ?

Câu trả lời là dùng context 
```go
func main() {
	now := time.Now()

	queue := 10
	ok := make(chan int, queue)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < queue; i++ {
		go reply(ctx, ok, i)
	}

	// ai nhanh thì thắng
	firstRepsonse := <-ok
	cancel()
	log.Println("since: ", time.Since(now), "pos", firstRepsonse)
}

func init() {
	rand.Seed(time.Now().Unix())
}

func reply(ctx context.Context, ch chan<- int, pos int) {
	rd := rand.Intn(3) + 1
	select {
	case <-time.After(time.Duration(rd) * time.Second):
		log.Println("run")
	case <-ctx.Done():
		log.Println("cancel")
	}

	ch <- pos
}
```
## Sử dụng channel để notification

**Đừng kể với ai nhé - 1-to-1 notification**

Một ngày đẹp trời, bạn được đi chơi riêng với cà rốt, có những điều từ lâu bạn đã muốn nói riêng với cô ấy. Hôm nay bạn mới có cơ hội vậy thì thắt dây an toàn và vào việc thôi.

```go
func main() {
	ch := make(chan string)
	go secret(ch)

	log.Println(<-ch)
}

func secret(ch chan<- string) {
	time.Sleep(time.Second)
	ch <- "this is secret"
}
```

**Gửi đám bạn thân - 1-to-n notification**

Bên trên bạn chỉ thổ lộ riêng với cà rốt, bạn nghĩ điều đó là bí mật riêng của hai người, nhưng mà biết đâu được, bạn đâu kiểm soát được nó. Bí mật chỉ là bí mật khi chỉ có một người biết thôi. Ví dụ cô ấy gửi tấm chân tình của bạn đến đám bạn thân chẳng hạn :v.

Có nhiều cách để thực hiện 1-to-n notification. Cách đơn giản nhất là close channel. Do tính chất: sau khi close channel thì receive value từ channel luôn luôn không bị block (receive zero value).

```go
func main() {
	ch := make(chan string)
	n := 10
	for i := 0; i < n; i++ {
		go notify(ch)
	}
	close(ch)
	time.Sleep(time.Second)
}

func notify(ch <-chan string) {
	<-ch
	fmt.Println("receive notify")
}
```

**Phản hồi về bí mật - n-to-1 notification**

Bí mật của bạn bị chia sẻ ra thì có lẽ nó cũng sẽ nhận được nhiều lời bàn tán. Cách dễ nhất để nhận nhiều lời bàn tán là dùng sync.WaitGroup. 

```go
func main() {
	n := 10
	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go notify(wg)
	}
	wg.Wait()
}

func notify(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("send notify")
}
```

**Sử dụng channels như Mutex Locks**

Channel với capacity bằng 1 có thể sử dụng như buffered channel. Tuy nhiên thì dùng cho vui thôi chứ thực tế chả ai dùng channel chỉ để implement lại cái mutex làm gì cả :3

```go
type empty struct{}

func main() {
	ch := make(chan empty, 1)

	cnt := 0
	wg := sync.WaitGroup{}
	n := 100
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				ch <- empty{}
				cnt++
				<-ch
			}
		}()
	}

	wg.Wait()
	log.Println("cnt: ", cnt)
}
```
Đoạn code trên dùng channel để đảm bảo synchronize cho biến cnt. Channel phải là buffered channel với cap = 1 để đảm bảo lần lock đầu tiên không bị block. Nếu không sẽ bị deadlock do ông nào cũng đòi khóa mà không ai chịu trả khóa.

## Sử dụng channels như Counting Semaphores

Buffered channel có thể sử dụng để implement semaphores. Counting semaphore có thể xem là multi-owner locks. Nếu cap của một channel là N thì nó có thể xem là một cái khóa cho phép N người sử dụng tại một thời điểm. Counting semaphores thường được sử dụng để giới hạn số lượng tài nguyên như số lượng concurrent request tối đa ...

Ví dụ khi bạn vào quán cafe, quán chỉ phục vụ một số lượng chỗ ngồi nhất định, hết chỗ thì ... mời bạn về. Dĩ nhiên rồi, chả nhẽ bạn lại đứng để uống cafe, mà có khi chỗ còn không có mà đứng ấy chứ. Thế làm sao để thông báo cho khách biết là đã hết chỗ rồi (ví dụ như khách đặt bàn online). Cách đơn giản là bạn phải duy trì một biến để đếm số lượng khách trong quán.

```go
type Seat int
type Bar chan Seat

func (bar Bar) ServeCustomer(c int) {
	log.Print("customer#", c, " enters the bar")
	seat := <-bar // receive value from bar
	log.Print("++ customer#", c, "drinks at seat#", seat)
	time.Sleep(time.Second * time.Duration(10+rand.Intn(6)))
	log.Print("-- customer#", c, " free seat#", seat)
	bar <- seat // release

}

func main() {
	rand.Seed(time.Now().UnixNano())

	bar24x7 := make(Bar, 10)
	for seatId := 0; seatId < cap(bar24x7); seatId++ {
		bar24x7 <- Seat(seatId)
	}

	for customerId := 0; ; customerId++ {
		log.Println("num goroutine: ", runtime.NumGoroutine())
		go bar24x7.ServeCustomer(customerId)
	}
}
```
Đoạn code trên đảm bảo chỉ có nhiều nhất 10 người được phục vụ tại một thời điểm. Mặc dù chỉ có nhiều nhất 10 người được phục vụ tại một thời điểm, tuy nhiên có nhiều hơn 10 customers ở trong hàng đợi để chờ phục vụ (> 10 goroutines), càng lâu thì số lượng goroutine này càng lớn. Từ đó sẽ bị tồn đọng và không bao giờ xử lý hết. Do tốc độ tạo mới nhiều hơn tốc độ consume.

Vậy làm sao để giới hạn số lượng goroutine có thể tạo ra? Ý tưởng là dùng một buffered channel với cap là số lượng channel lớn nhất có thể tồn tại tại một thời điểm.

```go
type Seat int
type Bar chan Seat

func (bar Bar) ServeCustomer(c int, seat Seat) {
	log.Print("++ customer#", c, "drinks at seat#", seat)
	time.Sleep(time.Second * time.Duration(2+rand.Intn(6)))
	log.Print("-- customer#", c, " free seat#", seat)
	bar <- seat // free seat and leave the bar
}

func main() {
	rand.Seed(time.Now().UnixNano())

	bar24x7 := make(Bar, 10)
	for seatId := 0; seatId < cap(bar24x7); seatId++ {
		bar24x7 <- Seat(seatId)
	}

	for customerId := 0; ; customerId++ {
		log.Println("num goroutine: ", runtime.NumGoroutine())
		// Need a seed to serve next customer
		seat := <-bar24x7
		go bar24x7.ServeCustomer(customerId, seat)
	}
}
```

Đoạn chương trình trên chỉ đảm bảo được có 10 goroutines tồn tại đồng thời. Tuy nhiên cần tạo rất nhiều goroutines trong quá trình hoạt động (mỗi khi free seat thì lại cần tạo một goroutine mới để xử lý).

Để tối ưu hơn, cần viết một chương trình đảm bảo chỉ có tối đa 10 goroutines tồn tại đồng thời và tối đa 10 goroutines được tạo ra. Đoạn code dưới đây thực hiện yêu cầu này.

```go
func main() {
	rand.Seed(time.Now().UnixNano())

	maxServe := 10
	consumers := make(chan int)
	for i := 0; i < maxServe; i++ {
		go serveCustomer(consumers)
	}

	for i := 0; ; i++ {
		log.Println("num goroutines:", runtime.NumGoroutine())
		consumers <- i
	}
}

func serveCustomer(consumers chan int) {
	for consumer := range consumers {
		log.Println("++ customer#", consumer, "drinks at the bar")
		time.Sleep(time.Second * time.Duration(2+rand.Intn(6)))
		log.Println("-- customer#", consumer, "leaves the bar")
	}
}
```
## Channel Encapsulated in Channel


## Try-Send và Try-Receive đến channel

Khi sử dụng select block với nhánh default và chỉ một nhánh case được gọi là try-send hoặc try-receive (tùy vào nhánh case triển khai ra sao). Try-send và try-receive không bao giờ block.

```go
func main() {
	type Book struct{ id int }
	bookshelf := make(chan Book, 3)

	for i := 0; i < cap(bookshelf)*2; i++ {
		select {
		case bookshelf <- Book{id: i}:
			fmt.Println("succeeded to put book")
		default:
			fmt.Println("failed to put book")
		}
	}

	for i := 0; i < cap(bookshelf)*2; i++ {
		select {
		case book := <-bookshelf:
			fmt.Println("succeeded to get book", book.id)
		default:
			fmt.Println("failed to get book")
		}
	}
}
```

**Check if a channel is closed without blocking the current goroutine**

Nếu có thể chắc chắn rằng không có values nào gửi đến channel nữa thì ta có thể sử dụng đoạn code sau (concurrently and safely) để kiểm tra xem một channel đã close hay chưa mà không cần block goroutine hiện tại.

```go
func isClose(c chan int) bool {
	select {
	case _, ok := <-c:
		return ok
	default:
	}

	return false
}
```
**Peak/burst limiting**

## Các cách khác để implement first-response win

Ta có thể sử dụng select mechanism (try-send) với buffered channel có capacity bằng 1 (ít nhất 1) để implement first-response-win. 

```go
func main() {
	rand.Seed(time.Now().UnixNano())
	c := make(chan int, 1)
	for i := 0; i < 5; i++ {
		go source(c)
	}
	rnd := <-c // only the first response win
	fmt.Println(rnd)
	select {}
}

func source(c chan<- int) {
	now := time.Now()
	defer func() {
		log.Println("end call to source after:", time.Since(now))
	}()

	fmt.Println("call to source")
	ra, rb := rand.Int(), rand.Intn(3)+1
	// Sleep 1s, 2s, 3s
	time.Sleep(time.Duration(rb) * time.Second)
	// try send
	select {
	case c <- ra:
	default:
	}
}
```

Phiên bản trên kết quả thì đúng nhưng sử dụng tài nguyên là thừa thãi. Do khi đã có first-response rồi thì không cần phải tiếp tục thực thi việc gọi đến source nữa.

```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan int, 1)
	for i := 0; i < 5; i++ {
		go source(ctx, c)
	}
	rnd := <-c // only the first response win
	cancel()
	log.Println("winner: ", rnd)
	time.Sleep(time.Second * 5)
	fmt.Println("number of goroutines:", runtime.NumGoroutine())
}

func source(ctx context.Context, c chan<- int) {
	now := time.Now()
	fn := func(ctx context.Context) <-chan int {
		rand.Seed(time.Now().UnixNano())

		fmt.Println("call to source")
		// try-send
		ch := make(chan int, 1)

		ra, rb := rand.Int(), rand.Intn(3)+1
		select {
		// simulate work load, we don't need to call when context is canceled
		case <-time.After(time.Duration(rb) * time.Second):
			defer func() {
				log.Println("end call to source after:", time.Since(now))
			}()
			ch <- ra
		case <-ctx.Done():
		}

		return ch
	}

	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	case v := <-fn(ctx):
		select {
		case c <- v:
		default:
		}
	}

}
```

## Rate Limiting

We can also use try-send to do rate limiting . In practice, rate-limit is often to avoid quota exceeding and resource exhaustion.

## Fan-in

Bạn cân viết một chương trình đi crawl data từ một source sau đó in ra console. 

```go
func main() {
	rand.Seed(time.Now().UnixNano())

	c := crawl()

	// forward to another source
	for i := 0; i < 5; i++ {
		log.Println(<-c)
	}

}

func crawl() chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("crawl %d", i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()

	return c
}
```
Yêu cầu mới là bạn phải crawl từ 2 source khác nhau, sau đó in ra console.

```go
func main() {
	rand.Seed(time.Now().UnixNano())

	c1 := crawl("crawl1", time.Second)
	c2 := crawl("crawl2", time.Second*5)

	// forward to another source
	for i := 0; i < 10; i++ {
		log.Println(<-c1)
		log.Println(<-c2)
	}

}
func crawl(id string, sleepTime time.Duration) chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("crawl for %s: %d", id, i)
			time.Sleep(sleepTime)
		}
	}()

	return c
}
```
Kết quả trả ra là một kết quả tuần tự crawl1,crawl2,crawl1,crawl2... crawl1 phải đợi crawl2 đi crawl xong thì mới được chạy.

```txt
2023/04/04 16:35:06 crawl for crawl1: 0
2023/04/04 16:35:06 crawl for crawl2: 0
2023/04/04 16:35:07 crawl for crawl1: 1
2023/04/04 16:35:11 crawl for crawl2: 1
2023/04/04 16:35:11 crawl for crawl1: 2
2023/04/04 16:35:16 crawl for crawl2: 2
2023/04/04 16:35:16 crawl for crawl1: 3
2023/04/04 16:35:21 crawl for crawl2: 3
2023/04/04 16:35:21 crawl for crawl1: 4
```

Sếp yêu cầu kết quả bạn crawl phải là realtime. Fan-in chính là giải pháp

```go
func main() {
	rand.Seed(time.Now().UnixNano())

	c1 := crawl("crawl1", time.Second)
	c2 := crawl("crawl2", time.Second*5)
	c := fanIn(c1, c2)

	// forward to another source
	for i := 0; i < 10; i++ {
		log.Println(<-c)
	}

}
func crawl(id string, sleepTime time.Duration) chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("crawl for %s: %d", id, i)
			time.Sleep(sleepTime)
		}
	}()

	return c
}

func fanIn(c ...chan string) <-chan string {
	fi := make(chan string)

	for _, v := range c {
		go func(v chan string) {
			for i := range v {
				fi <- i
			}
		}(v)
	}

	return fi
}
```
Kết quả đa số kết quả sẽ từ crawl1 vì crawl1 chạy nhanh hơn. 

```txt
2023/04/04 16:35:57 crawl for crawl1: 0
2023/04/04 16:35:57 crawl for crawl2: 0
2023/04/04 16:35:58 crawl for crawl1: 1
2023/04/04 16:35:59 crawl for crawl1: 2
2023/04/04 16:36:00 crawl for crawl1: 3
2023/04/04 16:36:01 crawl for crawl1: 4
2023/04/04 16:36:02 crawl for crawl2: 1
2023/04/04 16:36:02 crawl for crawl1: 5
2023/04/04 16:36:03 crawl for crawl1: 6
2023/04/04 16:36:04 crawl for crawl1: 7
```

Để sử dụng fan-in, luồng dữ liệu phải không phụ thuộc lẫn nhau


Đoạn code bên trên chưa graceful close channel và có goroutine leak. Chú ý phần này.
```go
func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithCancel(context.Background())

	c1 := crawl(ctx, "crawl1", time.Millisecond*1)
	c2 := crawl(ctx, "crawl2", time.Millisecond*2)
	c := fanIn(ctx, c1, c2)

	// forward to another source
	for i := 0; i < 2; i++ {
		log.Println(<-c)
	}
	cancel()
	// _ = cancel
	time.Sleep(time.Second * 1)
	log.Println("num goroutines:", runtime.NumGoroutine())

}
func crawl(ctx context.Context, id string, sleepTime time.Duration) chan string {
	c := make(chan string)
	go func() {
		defer func() {
			log.Println("exit crawl goroutine")
		}()
		for {
			select {
			case c <- fmt.Sprintf("crawl for %s", id):
			case <-ctx.Done():
				log.Printf("close %s\n", id)
				close(c)
				return
			}
		}
	}()

	return c
}

func fanIn(ctx context.Context, c ...chan string) <-chan string {
	fi := make(chan string)

	wg := sync.WaitGroup{}

	for _, v := range c {
		wg.Add(1)
		go func(v chan string) {
			defer func() {
				wg.Done()
				log.Println("exit fanIn goroutine")
			}()
			for i := range v {
				select {
				case fi <- i:
				case <-ctx.Done():
					return
				}
			}
		}(v)
	}

	go func() {
		defer func() {
			log.Println("exit wg goroutine")
		}()
		wg.Wait()
		close(fi)
	}()

	return fi
}
```
## The tee-channel


## Switches

Thao tác gửi hay nhận trên một nil channel đều bị block. Với tính chất này ta có thể kiểm soát được cơ chế selection trong select block.

```go
```

## Control code execution possibility weights

Ta có thể duplicate nhánh case trong select code block để tăng khả năng thực hiện đoạn code tương ứng.

```go
func main() {
	foo, bar := make(chan struct{}), make(chan struct{})
	close(foo)
	close(bar)
	x, y := 0.0, 0.0
	f := func() { x++ }
	g := func() { y++ }
	for i := 0; i < 10000; i++ {
		select {
		case <-foo:
			f()
		case <-foo:
			f()
		case <-bar:
			g()
		}
	}
	fmt.Println(x / y) // about 2
}
```

## Data Flow Manipulations