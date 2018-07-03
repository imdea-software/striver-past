PROGRAMS=( "AVGK" "STOCK" )
ARG2=( 1 10 50 100 )
MAXEVS=( 1000 10000 100000 1000000 )

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

