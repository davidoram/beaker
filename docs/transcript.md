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

- The core NATS messaging platform, which includes Jetstream persistent messaging
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

You will notice that the inventory microservice publishes changes as events to Jetstream.  This is an example of an 'event-driven' architecture. In our case when something interesting happens inside the inventory service (Low stock event), so that is published as a event into NATS, and we have configured NATS to capture those messages in an persistent message stream.  Why is that useful?  It allows consumers to connect to that stream and read the messages.  Because the messages are stored to disc, if the consumer goes offline, it can re-join later and catch up on all the messages it missed while it was offline.  Think of it like a holding bay for those messages, where they can be stored temporarily until the consumers are ready for them.   This style of architecture where we offer actions orient APIs which allow callers to do something actively, in conjunction with events that can be passively consumed allows for you to build applications on top of these services in a very flexible manner.

The last piece of the diagram shows NewRelic as our OpenTelemetry provider.  OpenTelemetry provides a standardized way for your app to collect metrics, logs, and traces about its behavior and performance. That information makes it easier to understand, debug, and monitor applications. If multiple systems all use open telemetry their information can be connected together.  Then you can choose the vendor that suits you best to view all that collected information. There are lots to choose from https://opentelemetry.io/ecosystem/vendors/ so you can choose based on price, or performance.  I've opted to go with NewRelic because I use it in my day job, so I know how to navigate around its interface.

I think of OpenTelemetry like the dashboard of your car.  You have all sorts of guages and warning lights that help you understand how teh engine (your application) is working. It will help us answer questions like, is the engine running, how fast id it workng, how much load is on the system, and is it about to break.

OK Thats enough about the high level architecture, lets focus in on the API that we are building:

## System Goals

If you scroll down on the same page to the "system Goals"  we will look at them one by one.

The goal is to build a basic inventory API that allows us to track quanity of different products.  It tells us that products are identified by a `product-sku`, and the three operation we have are to add, remove stock, and to query to see the current level.

The business rules section explains what constraints we want the system to maintain.

API endpoints explain each of the operations.

## Technical requirements 

Making the API accessaible via teh NATS and HTTP protocol is straightforward. Weve already discussed why is useful to give two options, but I didn't explain why you might want to encourage your callers to try the NATS protocol option. With NATS you get the following advantages:

- Lower latency & higher throughput. the NATS client keeps a persistent connection and uses a lighter weight protocol.
- Bidirectional comms. Once connected the caller can send and recieve data from the NATS server. It natively support streaming style operations.
- Resiliance. The NATS client will automatically route requests across the cluster giving us better availability and fault tolerance.
- Security. With a NATS persistent connection the callers authenticate once and stay secuerly authentcated rather then having to authenticate each call.


Our API requests and responses will use the JSON format. For those that don't know JSON is a human and machine readable format. https://en.wikipedia.org/wiki/JSON.  The Unix environment demonstrates how we can build powerful applications by using text oriented tools.  Many people fall into the trap of worrying that their data isn't compact enough over the wire and they are lossing a lot of performance, but that concern should be tempered wityh how easy it is to start using your system, constructing input payloads and understaning responses.  Using a text format like JSON means we can eyeball the data and check very easily if things look right.  If we were to use a binary format like protobuf then we need special tools to help us do that, and any extra layer of tools slows us down.  The other advantage to textual formats is that other tools can manipulate the data from our system.  There is a large list of JSON tools and systems available https://github.com/burningtree/awesome-json

OK so JSON is our data format, but we need to layer some structure over that. We need to validate all out API inputs and very strictly control our API outputs.  To help us we draw on another standard called JSON Schema https://json-schema.org. This allows us to define what our JSON data looks like and enforce those rules.   This makes our validation and conformance to easy to implement.  With each incoming request we get the JSON Schema library to check if its valid before proceeding.  Like JSON, JSON Schema is an open standard so its used in tools like OpenAPI to define API specifications.  By adopting a standard we are able to integrate with other systems and tools more easily.

Next on the list we have a requirement to implement tests to verify our system is working correctly. This is always a good idea, because it gives us confidence that when we make changes they are working correctly and we haven't broken something inadvertatly.   We will go through our test code in a later video.

Our last requirement is to capture telemetry using the open telemetry standard. We have already discussed what open telemetry is, and we wnat to use it so we cna oberve our application becahious at runtime.

## Standards and Organisations

If you scroll down in the page you will see some of the projects that this application uses, wither during development, or after its been deployed.  I added this so you can get a sense of what the projects are, and how they inter-relate.

Standards are useful because if they are widely adopted, and we make our applications follow them to then we get the benefits of joining a community.  For example I've used `git` and github as my source code control tool. Because gits a widely used tool, many organisations like github have standardised on its use, so I can take my git skills and use then inside my  github environment.

Lets look at the other tools, systems and standards.

This briongs us to the end of our video on our "High Level Architecture", I look forward to seeing you in the next video where we cover off the 'Development environment'.

Remember "Iron sharpens iron, and one man sharpens another.”. Hit the subscribe button if you wnat to be notified when the next video is out. See you next time.
