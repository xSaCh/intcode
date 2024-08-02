IN      x
SLT     x 1   y
JNZ     y :exit
ADD     x 0   z

loop:
ADD     z -1  z
JZ      z :print
MUL     x z x
JNZ     z :loop

print:
OUT     x
HALT    

exit:
OUT     1
HALT
HALT
