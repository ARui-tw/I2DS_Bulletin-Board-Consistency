# Project 2: Bulletin Board Consistency System
[GitHub Repo](https://github.com/ARui-tw/I2DS_Bulletin-Board-Consistency)
## Team members and contributions
- Group: 6
- Team members:
	- Eric Chen (chen8332@umn.edu)
	- Rohan Shanbhag (shanb020@umn.edu)

### Member contributions:
We opted for a pair programming approach towards the coding/development process by working together through Discord, meeting generally every night for an hour or two over the course of the project. Both team members shared equal responsibility for overseeing the code written, the design decisions/documentation, and the test cases. The two of us were working together on call to complete the Post(), Read(), Choose(), Reply() functionality.


## Project Build and Compilation Instructions
In order to run the project:
### Update Go
Since the version on the lab machines is outdated, we need to update the Go version to the latest version. I've written a script to do this for you. Just run the following command:
```sh
. build-go.sh
```
**NOTE**: This script will modify the $PATH variable to point to the new Go installation. Once you exit the shell, you will need to go to the project's folder and run the following command again to update the $PATH variable. (I don't want to mess up with the $PATH variable in your shell profile, but you can add the following line to your shell profile if you want to make it permanent.)
```sh
export PATH=$PWD/BuildGo/go/bin:$PATH
```
### Run the server: 
To configure and run the server, first edit the `server_config.json` (the naming of the file is up to you). There are two examples file in the repo, `server_config.json` is for primary-backup protocol and the `server_config_q.json` is for the quorum consistency. Then run the following command:
```sh
go run main.go -config_file server_config.json
```
Wail till the message "All servers started" is printed, then you can start the client.

### Run the client: 
Just simple run the following command and client will join random server, note that the config file will have to be the same as the one used to run the server:
```sh
go run BBC_client/main.go -config_file server_config.json
```

## Design Document
We decided upon the design of our Bulletin Board Consistency System for Project 2 between the two of us, and some specific design decisions we made were as follows:
* The server's port is predefined in the config file.
* The client has the option to Post, Read, Choose or Reply to articles. Upon posting an article, the client is shown a hierarchical list of output, with any posts in the outermost level of the hierarchy, and any replies indented further inwards.
* The articles are presented to the client after the Post or Read requests in time order (along with their ID).
* When an article is added the first few words are shown along with its ID
* Our servers can handle subsequent read/post requests, allowing a client connected to one server to read, and for another client to post, while outputting the expected changes to each client
* Our server can handle concurrent requests from clients (to post/read/reply), and our elected method of handling concurrent requests is to let the server allocate resources on a first-come-first-served basis. This does not prioritize a certain request, and instead requests to post/read/reply are given equal priority.
* We have chosen to pre-define a set of 10 servers that we have on standby for use with our Bulletin Board system.


## Simple Analysis
### Primary-Backup Protocol (Sequential Consistency & Read-your-Write consistency)
As part of our Primary-Backup Protocol implementation, we satisfied both Sequential Consistency and Read-your-Write consistency. We can start up as many server as we want, but in this way, it will take more time when reading (Post, Reply). However, since every time the client want to write, we will push the change to each back-up server, the read time will be faster.
### Quorum Consistency
As part of our quorum consistency implementation, we had chosen to have 8 servers running, where 4 servers existed as part of the read quorum and 5 servers existed as part of the write quorum – one server was both a read and write server. We changed these numbers around (as indicated by a few of the different cases we outlined below) in order to generate a performance graph to identify the average throughput of the read and write operations, and also to attempt to detail an ideal ratio of Nr to Nw given our implementation of Quorum consistency.

- Case 1:
	- N = 8
	- Nr = 4
	- Nw = 5
- Case 2:
	- N = 6
	- Nr = 3
	- Nw = 4
- Case 3:
	- N = 10
	- Nr = 1
	- Nw = 10
- Case 4:
	- N = 3
	- Nr = 2
	- Nw = 2

### Analysis:
For the quorum consistency, it will speed up the writing time since instead of updating the data to all servers, it will only have to update it to Nr servers. However, the tradeoff is that the reading speed will be slower since it will have to check with Nr servers to make sure that the data is consistent. But the advantage of this protocol is that the ratio is flexible, so we can change the ratio to make the reading/writing time faster/slower. For example, in case 3, the writing speed will be same as the primary-backup protocol.

| Case    | Writing Speed Ranking | Reading Speed Ranking  | Fault Tolerance Ranking  |
| ------- |:---------------------:|:----------------------:|:-:|
| Case 1  | 3                     | 4                      | 2 |
| Case 2  | 2                     | 3                      | 3 |
| Case 3  | 4                     | 1                      | 1 |
| Case 4  | 1                     | 2                      | 4 |

## Test Cases
In the test cases outlined below, whenever/wherever a user enters an article, our program at that time asks our user to provide a file path. The Readme/Test Cases outlined below simply show the contents of the article that are being posted instead in order to aid with its overall readability. Additionally, since these are for readability/explanation purposes, the actual test data (articles that are being posted/read) may be different, but the output will be in a similar format/manner to that listed below.
### Test Case #1: Normal Condition
This test case tests sequential consistency, and ensures that the Bulletin Board System is consistent and runs as intended.
#### How to run:
* Start up the server
	```json
	{
		"type": "PBP",
		"primary": 50051,
		"child" :[
			50052,
			50053,
			50054
		]
	}
	```
* Client 1:
	* Join random server
	* \<user enters\> "1"
	* \<user enters\> "Sample/1.txt"
	* \<user enters\> "2"
	* Server displays hierarchy of articles, showing:
	```
		0: Series 05 Episode 04 – The Wiggly Finger Catalyst
	```
	* \<user enters\> "4"
	* \<user enters\> "0"
	* \<user enters\> "Sample/2.txt"
	* \<user enters\> "q"
* Client 2:
	* Join random server
	* \<user enters\> "2"
	* Terminal displays hierarchy of articles, showing:
		```
		0: Series 05 Episode 04 – The Wiggly Finger Catalyst
		1: └─The One Where Monica Gets a New Roommate (The Pilot-The Uncut Version)
		```
	* \<user enters\> "1"
	* \<user enters\> "Sample/3.txt"
	* \<user enters\> "q"
* Client 3:
	* Join random server
	* \<user enters\> "2"
	* Terminal displays hierarchy of articles, showing:
	```
	0: Series 05 Episode 04 – The Wiggly Finger Catalyst
	1: └─The One Where Monica Gets a New Roommate (The Pilot-The Uncut Version)
	2: The Junior Mints
	```
	* \<user enters\> "4"
	* \<user enters\> "1"
	* \<user enters\> "Sample/4.txt"
	* \<user enters\> "2"
	* Terminal displays hierarchy of articles, showing:
	```
	0: Series 05 Episode 04 – The Wiggly Finger Catalyst
	1: └─The One Where Monica Gets a New Roommate (The Pilot-The Uncut Version)
	2: │ └─50% Off
	3: The Junior Mints
	```
	* \<user enters\> "q"


### Test Case #2: Failure Condition
An error message will be printed in the negative case that a connected client tries to choose an article that does not exist, and the user will be notified to only attempt a reply when other articles are posted/readable.
#### How to run:
* Start up the server
* Client 1:
	* \<user enters\> "3"
	* Server prints out error message: "No article to choose, please read first", and lets user retry the system


### Test Case #3: Edge Condition 
Testing the Quorum consistency, we will have 4 servers as part of the read quorum (Nr) and 5 servers as part of the write quorum (Nw). We can have any number of clients join, but for the sake of our test case, let us say there are three clients
#### How to run:
* Start up the server
	```json
	{
		"type": "quorum",
		"primary": 50051,
		"child" :[
			50052,
			50053,
			50054,
			50055,
			50056,
			50057,
			50058
		],
		"NW": 5,
		"NR": 4
	}
	```
* Client 1:
	* \<user enters\> "1"
	* \<user enters\> "Sample/1.txt"
	* \<user enters\> "q"
* Client 2:
	* \<user enters\> "2"
	* Terminal displays hierarchy of articles, with the first served post request appearing with ID 1. Let us assume Client 1’s request goes through first, then the system would show:
	```
	0: Series 05 Episode 04 – The Wiggly Finger Catalyst
	```
	* \<user enters\> "q"
* Client 3:
	* \<user enters\> "2"
	* Terminal displays hierarchy of articles, with the first served post request appearing with ID 1. Let us assume Client 1’s request goes through first, then the system would show:
	```
	0: Series 05 Episode 04 – The Wiggly Finger Catalyst
	```
	* \<user enters\> "1"
	* \<user enters\> "Sample/2.txt"
	* \<user enters\> "2"
	* Terminal displays hierarchy of articles, with the first served post request appearing with ID 1. Let us assume Client 1’s request goes through first, then the system would show:
	```
	0: Series 05 Episode 04 – The Wiggly Finger Catalyst
	1: The One Where Monica Gets a New Roommate (The Pilot-The Uncut Version)
	```
	* \<user enters\> "q"
* Client 4:
	* \<user enters\> "2"
	* Terminal displays hierarchy of articles, with the first served post request appearing with ID 1. Let us assume Client 1’s request goes through first, then the system would show:
	```
	0: Series 05 Episode 04 – The Wiggly Finger Catalyst
	1: The One Where Monica Gets a New Roommate (The Pilot-The Uncut Version)
	```
	* \<user enters\> "4"
	* \<user enters\> "0"
	* \<user enters\> "Sample/3.txt"
	* \<user enters\> "q"
* Client 5:
	* \<user enters\> "2"
	* Server displays hierarchy of articles, showing:
	```
	0: Series 05 Episode 04 – The Wiggly Finger Catalyst
	1: └─The Junior Mints
	2: The One Where Monica Gets a New Roommate (The Pilot-The Uncut Version)
	```
	* \<user enters\> "q"
* Client 6:
	* \<user enters\> "2"
	* \<user enters\> "3"
	* \<user enters\> "1"
	* Server displays the output of the second article
	* \<user enters\> "q"

## Pledge
No-one sought out any on-line solutions, e.g. GitHub for portions of this lab

Signed:
Rohan Shanbhag
Eric Chen