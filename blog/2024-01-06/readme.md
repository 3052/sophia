# jan 3 2024

start with this:

~~~
youtube -b BCRhBaFqtf0 -vc avc1
~~~

fragment it:

~~~
ffmpeg -i BCRhBaFqtf0.mp4 -c copy -frag_size 9999 -movflags dash frag.mp4
~~~