# id3
Ross Quinlan's ID3 decision tree in pure Go.

Run `make` to compile and run.

All commits should be checked by Github Actions and builds are available [here](https://github.com/nwsnate/id3/actions).
Actions is in beta, and isn't working properly all the time.

### Writeup

Writing id3 in Go was a challenge, given its the first large project that I've done in Go. Along the way I discovered many of the intricacies of Go, as well as developing my own taste for how languages should work. Working in a lower level language like Go gave me a new insight into a part of computing that I previously had little experience with.

I had a lot of fun learning a new language, but for the first project of the year and on a topic that is new to me, maybe another language would have been a better choice. I started out the project ahead of others, and finished some of the basic functions very quickly. Then I realized Go's default list implementation isn't fully developed, and had to write a few functions to interact with Go's arrays. This along with other issues with the Go compiler, caused a bit of a slow down in progress.

All in all, I would definitely use Go again. I think that it shows a lot of promise in the future, due to some big projects (Kubernetes, Promethus, Docker) by big name companies. Similar to Python, where Go really shines is 3rd party libraries. Support for a vast variety of functionality has been added to Go in a very short time since it's release. I think the problems I faced were a result of lack of experience in low level languages, and my choice to use Go's default arrays instead of one of the standard library packages such as `collection/list`. I look forward to learning even more about the language, and I think that id3 not only helped me gain knowledge about the algorithm itself, but insight into a new language that I will definitely continue to use in the future.

Below is a table of some trials I ran and how the algorithm performed.


What datasets did you use?  How many training and testing examples did you have?  What accuracy did you achieve?  Did you try different proportions of training and testing examples?  What effect did that have on your testing set accuracy?  How big were the trees that you got (depth, number of nodes)?

| Dataset | Training | Accuracy |
| ------- | -------- | -------- |
| tennis.txt   |                     |          |
| tumor.txt    |                     |          |
| mushroom.txt |                     |          |

