---
title: "Một vài cách sử dụng channel trong golang"
date: 2023-03-03
authors: ["ngoctd"]
draft: false
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

Chú ý là phải dùng buffered channel để tránh ăn send block của channel nhé. Đoạn code trên còn vấn đề là chẳng hạn cà rốt chỉ đăng tin cho người nào đó xem thôi, người đó reply rồi thì cà rốt thêm một tin mới là đã tìm được người cần tìm rồi. Thế là bạn cũng không cần reply lại nữa. Vậy xử lý thế nào nhỉ?

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

**Sử dụng channels như Counting Semaphores**