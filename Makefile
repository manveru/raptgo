include $(GOROOT)/src/Make.$(GOARCH)

TARG=raptgo
GOFILES=\
	blitz.go \
	enemies.go \
	enemy.go \
	flak.go \
	goodie.go \
	level.go \
	micro_missile.go \
	mine.go \
	missile.go \
	raptgo.go \
	raptor.go \
	reaver.go \
	shape.go \
	surface.go \
	thor.go \
	tracking_flak.go \
	tsunami.go \
	vector.go \
	weapon.go \
	odin.go \
	mauler.go \
	deathray.go \
	eclipse.go

raptgo: all
	$(GC) -I ./ -I /home/manveru/github/banthar/Go-SDL/ -I /home/manveru/github/banthar/Go-SDL/sdl ${GOFILES}
	$(LD) -L ./ -o raptgo _go_.$(O)


include $(GOROOT)/src/Make.pkg
