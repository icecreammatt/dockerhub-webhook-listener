#!/bin/sh
echo "DEPLOYING|$1|" >> deployed.txt
echo `date` >> deployed.txt
echo Deployed `date`
#touch `date`-testFile.txt

#docker pull icecreammatt/lookup
#docker run --name lookup -d -p 5003:5000 icecreammatt/lookup
