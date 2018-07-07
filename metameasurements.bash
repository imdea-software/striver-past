PROGRAMS=( "STOCK" "AVGK" )
ARG2=( 10 50 100 500 )
# ARG2=( 200 210 220 230 240 250 260 270 280 290 300 310 320 330 340 350 360 370 380 390 400 410 420 430 440 450 460 470 480 490 500 510 520 530 540 550 560 570 580 590 600 )
MAXEVS=( 500000 5000000 50000000 )
# MAXEVS=( 1000000 )

progi=0
while [ "x${PROGRAMS[progi]}" != "x" ]
do
    prog=${PROGRAMS[progi]}
    progi=$(( $progi +1 ))
    arg2i=0
    while [ "x${ARG2[arg2i]}" != "x" ]
    do
        arg2=${ARG2[arg2i]}
        arg2i=$(( $arg2i +1 ))
        maxevsi=0
        while [ "x${MAXEVS[maxevsi]}" != "x" ]
        do
            maxevs=${MAXEVS[maxevsi]}
            maxevsi=$(( $maxevsi +1 ))
            sh measurements.sh $prog $arg2 $maxevs
        done
    done
done

