.PHONY: build clean
build:
	mkdir -p build lib
	cd build && cmake -S .. -B . && cmake --build .
	cp build/libpqdevkit.a lib/libpqdevkit.a

clean:
	rm -rf build lib