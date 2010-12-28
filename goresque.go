package main

import (
	redis "clarity/drivers/redis"
	"fmt"
	"os"
	"strconv"
)

type Resque struct {
	Server string
	Port   int
	Db     int
	Queues []Queue
	Workers []Worker
	client *redis.Client
}

type Queue struct {
	Id     int
	Name   string
	client *redis.Client
}

type Worker struct {
	Id	int
	Name string
	client *redis.Client
}




func (self *Queue) size() (int, os.Error) {
	key := fmt.Sprintf("resque:queue:%s", self.Name)
	return self.client.Llen(key)
}

func (self *Resque) getStat(name string) (int, os.Error) {
	key := fmt.Sprintf("resque:stat:%s", name)
	val, err := self.client.Get(key)
	strval := string(val)
	intval, _ := strconv.Atoi(strval)
	return intval, err
}

func (self *Resque) getWorkers()([]Worker){
	workers, err := self.client.Smembers("resque:workers")
	if err != nil {
		fmt.Println(err)
	}
	var w Worker
	qs := make([]Worker, 1000)
	for i, val := range workers {
		w = Worker{Id: i, Name: string(val)}
		w.client = self.client
		qs[i] = w
	}
	self.Workers = qs
	return self.Workers[0:len(workers)]
}


func (self *Resque) getQueues() ([]Queue) {
	members, err := self.client.Smembers("resque:queues")
	if err != nil {
		fmt.Println(err)
	}
	var q Queue
	qs := make([]Queue, 100)
	for i, val := range members {
		q = Queue{Id: i, Name: string(val)}
		q.client = self.client
		qs[i] = q
	}
	self.Queues = qs
	return self.Queues[0:len(members)]
}

func NewResque(server string, port int, db int) (resque *Resque) {
	resque = new(Resque)
	resque.Server = server
	resque.Port = port
	resque.Db = db
	resque.client = new(redis.Client)
	resque.Queues = make([]Queue, 0)
	address := fmt.Sprintf("%s:%d", resque.Server, resque.Port)
	resque.client.Addr = address
	return resque
}

func main() {

	r := NewResque("web6", 6379, 0)

	queues := r.getQueues()

	for _, v := range queues {
		i, _ := v.size()
		fmt.Println("[", v.Id, "]", string(v.Name), "(", i, " items)")

	}
	workers := r.getWorkers()

	for _, v := range workers {
		fmt.Println("[", v.Id, "]", string(v.Name))

	}

	processed, _ := r.getStat("processed")
	failed, _ := r.getStat("failed")
	pending, _ := r.getStat("pending")


	fmt.Println("Processed:", processed)

	fmt.Println("Failed:",failed)
	fmt.Println("Pending:",pending)
	fmt.Println("Queues:",len(queues))
	fmt.Println("Working:",0)
	

}