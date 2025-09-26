# Transacript

# Espisode 1 - Introduction and overview

Hi and welcome to my series on "Production grade system development"

My name is Dave Oram, and in this series I'm going to walk through the elements of developing a production grade microservice API system.  Instead my aim is to be more holistic and descibre why I'm making some decisions. Yes we will look at all of the code, but we will also cover important issues around design, developer environment setup, production grade runtime monitoring and system integration concerns.

For this reason the series is useful to CTOs, developers, architects, product managers, test engineers or anyone interested in the whole development process.

Each video looks at a different aspect of the system, and by the end of it we will have looked ate every file in the project, so you should be able to get a sense for how the system works.  There will be lots of discussion and talk about what options were considered and what I chose and why.  These are just my opions and what works for me. I've been developing systems for over 35 years, and what I'm presenting today models the systems I've set up in my day job. All these decisions involve personal choice and of course you should feel free to pick and choose the ideas and techniques that work for you.  The idea is to help develop some questions to ask yourself and your colleagues when making deicisons, so that you can feel confident in making good decisions for you.

OK, with that preambe out of the way lets get some basic housekeeping out of the way and get started.

All of the code that we will be walking through is publicaly available on github at https://github.com/davidoram/beaker. In a later video we will talk about the online development environment we will use, but for now you can just navigate to that page in the browser and take a look.  The project is MIT licensed.  This means you are free to borrow any or all of the code for your own use.  

OK, I've covered the introduction to this series, I look forward to seeing you in the next video where we cover off the 'High level architecture'.

Remember "Iron sharpens iron, and one man sharpens another.”. Hit the subscribe button if you wnat to be notified when the next video is out. See you next time.

# Episode 2

Hi and welcome to my series on "Production grade system development". My name is Dave Oram and I'll be your gude as we look at episode 2 covering our "high level architectural goals".

If you want to learn more about this video series, I encourage you take a look at that, its only a few minutes long.

OK, today we will be covering off the high level architecture. If you want to follow along point your browser at https://github.com/davidoram/beaker, or just watch as we walk through the our high level arhcitecture.

Starting at the README.md file which is displayed when you navigate to the github project.  There are a few statements in this file which I want to talk about.

The overview talks about a "production ready microservice API". What do I mean by that? Well there are a million tutorials on the web about creating APIs, but whats the difference between demo code and production code. I think of it as the difference between say a classic car, and a delivery truck. The classic car might look a lot cooler, but it needs a lot of TLC to keep it running well. Every few miles you need to top up the oil because its got a leak somewhere, you can only hold two people in the classic car because its a fancy convertable. On the other hand the truck is working day in day out. It needs to take heavy loads when necessary, and it needs to be super reliable.  There is a team of mechanics ready to jump in and fix the truck when it breaks because time is money with that expensive truck.  

 In the context of an API we might need to consider how that API fits into its environment. We want that API to be able to scale up when necessary and handle great load.  While that API is running we want tools to evaluate its peformance and health so we can be confident that its working correctly.  With a little bit of planning we can create a API that allows us to scale as our API gets more popular. In a future video we will be covering how we capture Telemetry information which provides real time diagnostic informtaion about our system as its running.

 In the next section I talk about three technology cornerstones that we ise to build our system. The go programming language, the NATS messaging system and Postgres database.

 Why do I choose to program in go.  Well I've had many years experince with go so I'm well placed to compare it with other languages I've used in the past.  For me go is an excellent language for modern enterprise system development. There are many cool features built into go which you cna read about elsewhere, but I'm going to raise a few that affect productivity over the long term:
 
 - For starers, having a static typing eliminiates a whole class of run-time bugs, which for me means increased reliability. 
 - The tooling is second to none, and includes source code formatting, unit testing, linting, security vulnerability checking is built in.  This means teams can spend less time figuring out for example how to format code, or write unit tests, of check for security vulenabilities - those decisions are made for you so you can concentrate on your business problems.
 - Go version compatability guarante. A new version of go is released every 6 months, you upgrade by recompiling and typically what that means is you get increased performance for free. Other languages evolve their syntax and libraries quickly, but it means you are constantly having to update code for no benefit apart from keeping current. And you need to keep your code current because production systems need to use the latest versions to stay ahead of security vulnerabilities.

 OK lets click on the link to the "High Level architecture".

## High level Architecture

Lets look at all the pieces on the diagram in front of us.

OK we are building a microservice architecture.  What that means is we build small components, which helps you as a developer keep each component small and easy to understand.  We have boundaries between our microservices and if we keep the boundary interactions the same, then we can change the internal structure of any component without affecting the others.  As systems get larger that can be a great help in managing complexity. In this project we have a single micro-service for managing inventory.  If in the future we wanted to manage customers, we would make that a separate micro-service with its own application executable and postgres database. 

For the purposes of this exercise we will be running our microservice inside our codespace development environment. We will talk more about our development environment in a later video.

Our microservice stores its persistent data in the Postgres relational database system. Our database design and tooling will be covered in a separate video.

Our microservice code needs a way to communicate with the outside world, and we do that through NATS. NATS is an open source messaging system, and among its many capabilities it allows us to build API services which work on a request response model. To start with, all interactions with NATS are secured and require callers to be authenticated and authorized by NATS. To call an API a request is sent to NATS, that request is routed to the correct microservice which receives the request and it replys with a response that NATS routes back to the caller.   This is how regular HTTP APIs work as well, its just that the transport layer isn't HTTP, instead our microservice uses the NATS protocol which is optimized for this kind of message exchange. What that means is that to call the API, by default the caller needs to use a NATS client rather than an HTTP client.  If you are an experienced developer at this point you may be asking yourself, hey HTTP is the lingua franca protocol for APIs on the web, so why would I want to use NATS and have my clients have to use a different protocol.

Well this is where Synadia come into the picture. They are the commercial company behind NATS and they offer a hosted NATS service caled Synadia Cloud. A hosted service means that you pay Synadia and they will run NATS on your behalf.  Lets take a look at what a Synadia Cloud gives us:

- The core NATS messaging platform, including Jetstream messaging storage
- A Web UI Console for management
- HTTP Gateway
- Connectors
- Custom workloads

Back to the question of offering an HTTP API.  Synadia Cloud includes the 'HTTP Gateway'. Gateway is a common term that means a system or devixe that acts as a bridge beyween two different protocols.  In this case the HTTP Gateway converts incoming HTTP requests to NATS. Think of it as a 'translator' between the two protocols.

Using the HTTP Gateway doesn't require any extra work on your behalf, so its nice that you get dual protocol support for free.  If you are an enterprise customer, it's a killer feature because its often easier for customers to integrate with HTTP APIs rather than NATS because its so familiar to people.

Of course there are always tradeoffs that have to be made, and with the HTTP Gateway there are several to consider:

- You have no control of the URL, or the format of the API endpoint. Its always a  `PUT https://api.ngs.global/nats/subjects/{subject}` 
- You have no control, and no visibility over errors returned by the HTTP Gateway component. For example if your microservice goes offline then callers through the HTTP Gateway will get a `500 no responders` message.

OK, so like all the decisions we make across this project we weigh up the pros and cons and decide if this works for us.

For me, one of the main drawcards for using Synadia Cloud is that they provide, update and monitor all the infrastructe and settings needed to provide secure access to NATS.  At first glance it might seem easy, but here are just some of the things that they are doing for you when you use their service:

- They obtain and manage SSL certificates used by the TLS security layers to client connections secure.
- They provide a cluster of connection points to NATS, which is geographically distributed, and also multi-cloud. This provides redundancy and reliability
- They manage the Jetstream persistant message storage for us.
- Their teams are monitoring for security vulnerabilities, constantly patching the software stack, and also regularly upgrading NATS

Of course not everyone needs this kind of service, but it can save your team a lot of time and effort if these things are already convered off. You can focus on building a great product.

You will notice that the inventory microservice publishes changes as events to Jetstream.  This is an example of an 'event-driven' architecture. In our case when something intersting happens inside the inventory service