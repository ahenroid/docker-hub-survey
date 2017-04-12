GO := go

EXE := pull-stats base-stats median-stats

all: $(EXE)

clean:
	$(RM) $(EXE)

%: %.go
	$(GO) build $<

.PHONY: all clean
