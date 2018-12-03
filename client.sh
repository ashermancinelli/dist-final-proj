while read line
do
	if [ "$line" == "exit" ]; then 
		echo "Recieved kill signal..."
		break
	else
		echo "$line"
	fi
done < <((echo "Welcome!") | nc $1 $2)


