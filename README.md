# Created a Music Visualiser

## Supports only wav file for now (Other file support to be added)

* Uses DFT(Discreet Fourier Transformation) to convert amplitude-time graph to frequency-time graph.(FFT upgradation can be done here)
* Displays data on a webpage by fetching data using a websocket connection continuosly fetch and update data.
* The data displayed couldn't be fetched on a constant basis as it will lead to data scarcity since DFT takes time and the data fetching is faster than that.
* Implemented a queue which will fetch data for a while till there is enough data to share to the webpage to avoid data scarcity.


## DATA FLOW
1. Webpage makes a get Request to backend server 
2. Backend server starts sending data to queue and when enough data is sent, sends reply back to webpage.
3. Webpage on recieving reply sends websocket request to queue server, which starts sending data to webpage.
4. Webpage renders the data and ta-da.
