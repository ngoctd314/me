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
