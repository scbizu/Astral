package tl

type Stash chan PMatch

var (
	stashMatchChan Stash
)

func GetStashChan() Stash {
	if stashMatchChan == nil {
		stashMatchChan = make(chan PMatch)
	}
	return stashMatchChan
}

func (s Stash) Put(m PMatch) {
	s <- m
}

func (s Stash) Run() {
	for {
		select {
		case <-s:
			// TODO: terminated the match
		}
	}
}
