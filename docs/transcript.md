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

Hi and welcome to my series on "Production grade system development". My name is Dave Oram and I'll be your gude as we todays espisode which covers our "high level architectural goals".

If you want to learn more about this video series, I encourage you take a look at the previous video which provides a brief introduction to all teh things that we will be covering.

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


# Episode 3

Hi and welcome to my series on "Production grade system development". My name is Dave Oram and I'll be your gude as we todays espisode which covers our "development environment".

If you are new to this video series video series, I encourage you take go back and listen to previous videos as they cover some context to what we are covering today.

OK, today we will be covering off the "development environment". If you want to follow along point your browser at https://github.com/davidoram/beaker, otherwise you can just watch me cover all the steps.

Starting at the README.md file which is displayed when you navigate to the github project. Then click on the "Development environment" link.

What is a developoment enviroment?  Is where your developers, testers, and infrastructure engineers write, test and debug code.

A good development environment contains all the tools, and resources that you need to build, run, debug and test the application safely without affecting anyone else.

Ideally that means that each person can work concurrently on a separate feature or bug fix without impacting others.   Isolation of development environments helps person work at their own individual speed without getting blocked by someone else.  I've worked in many places in the past where there was always some part of the system where developers or testerd bumped into each other. Often this happens at the data layer, for example multiple developers / testers work against a shared database.  This slows down development because of issues like:

- Each tester must take care not to overwrite or change someone elses test data.
- Developers struggle to co-ordinate database schema changes when they impact others.
- Performance testing becomes impacted by other peoples workloads so accurate measurements become harder to obtain.

All this leads to increased friction and decreased productivity, so our goal is to create  isolated development environments.

It takes a lot of effort to create a development environment. Consider what might happen when a new person is joining your team and they needs to get their development environmnet setup.

This might involve:
- Operating system updates
- Containerization software installs eg: Docker, Podman or Kubernetes
- Code editors or Integrated Development environments, eg: VS Code, or Vim
- Debuggers
- Source code control systems eg: git
- Compilers: eg: C#, python, or go
- Databases eg: Postgres, or MySql with associated client tools
- Specialized tools used inside your development flow, maybe to generate code, perform liniting, or produce documentation.
- SaaS credentials, eg: to access your GitHub acount, JIRA system for ticketing, and an cloud provider like AWS or Azure, etc.

Once you have done that then you are ready you might be just about ready to checkout the source code for your application and start working.

The problem with all this setup is that you end up with these problems:

- Complex Setup: New team members face a long, error-prone process to install and configure all required tools, dependencies, and credentials.
- Maintenance Burden: Keeping environments up-to-date (e.g., upgrading Postgres) is difficult and often inconsistent across the team.
- Troubleshooting Difficulties: Diagnosing and resolving environment-specific issues is hard, especially when problems are unique to one developer’s setup.
- Lack of Isolation: Sharing environments or resources (like databases) can lead to conflicts, data corruption, and reduced productivity.
- Experimentation Risk: Trying new versions of tools or languages can disrupt existing setups, making it hard to safely experiment without breaking things.


In short, Setting up and maintaining consistent, reliable, and isolated development environments is challenging, leading to wasted time, hard-to-diagnose issues, and friction for both onboarding and day-to-day work.

So this is why many development teams (myself included) have moved to online "Cloud Based development environments".

The idea behind this is that the team builds a standardised development environment typically using containerization technoloy, and then when a developer needs an environment they spin one up, use it for the duration of a piece of work, and then discard it once done.

Each cloud based development environment instance runs in the cloud, but looks ans acts just like a normal development environment that you used to run on your laptop.

There are many vendors that offer this service including
- Github codespaces https://github.com/features/codespaces
- Google Cloud Workstations
as well as open source soultions like coder https://coder.com

I've used a couple of these but I'm going to focus on GitHub codespaces because its the environment I used in my day job, and its the one I'm most familiar with.

If you scroll down the page you can see a diagram showing the GitHub codespace environment that we will be running.

We will be running either vscode or a modern browser on our desktop. Chrome or Safari should work fine.

When we start the codespace, a virtual machine is started somewhere in Githubs infrastructure, and on that machine it has the following:
- First of all it has git, and all of our source code of the beaker project checked out.
- It has a bunch of tools installed like go compilers, and command line tools that we need to build, run and test the system
- It has a full docker environment, and inside that docker environment we run Postgres. This is going to hold our test and production databases.

You will notice on the left of the diagram you see that our codespace connects to some other services on the internet:
- Synadia Cloid - for our NATS service
- New Relic - for our Open Telementry data
- Github obviously for recording our source code changes.
- Not shown on the diagram, but codespaces has automatic port forwarding so for example you can run a website on your codespace and view it in your local browser. 


So what are the advantages of our coder environment:
- Your developer only needs a browser and an internet connection to get up and running. This makes onboarding a one click operation.
- IT managers might issue Macbook Airs ($1,800) vs MacBook Pros ($2,819). A saving of ~ 36%
- Each developers environments is completely consistent on startup - because we script it
- Development environment setup is scripted so you can experiment and make changes with the knowledge you can back them out if they don't work.  You can not only change the software, but you can change the virtual hardware, increasing the number of CPUs or memory. 
- Enterprises can control access and costs
- Developers can run multiple independent codespaces concurrently. They autoamtically 'pause' and shutdown so you just pay for what you use.

But there are some disadvantages.
- Will developers accept a standardised set of tools. This is especially relenant to the IDE, so codespaces works best on VSCode and has beta support for JetBrains IDE, but what if your devs use VIM?
- It costs money. Personal accounts get a quota of free usage and its very afforable afterwards. For example I've been using codespaces for this project and my bill last month was 0.49 cents! Enterprise plans of course cost real money so that needs to be factored into the equasion.


I've been using codespaces at my day job for over 9 months now and for me it works extremely well. You'll need to make your own judgements about how it works for you.

OK, lets dive into a beaker projects codespace settings and examine  step by step how its put togther.

Click on the link to open the `dev-container.json` file. If you have checked out the project already it lives in thge `.devcontainer` directory in the root of the project.

Ok we start off with the `name` of our project which is informational only, so I called it 'Beaker project dev container'

Next the `workspaceFolder` shows where all the application project files will live.

The `Features` section pulls in pre-packaged tools.  There are lots of them to choose from and I have added some links into the documentation.  This is the simplest way to start adding tools into our development environment. Ive added the `go` development environment and something called `docker-in-docker`.
The `go` feature adds the go programming language tools. `docker-in-docker` runs the docker daemon process inside docker, and allows our codespace environment to build and run docker containers.

The next section is `forwardPorts`. A forwarded port moves network traffic from one computer to another. In this context we are telling the codespace to forward network traffic from the codespace virtual machine to the developers laptop. Port 5432 is used to query the Postgres relational db.  So this means I can run queries from my laptop they will be re-directed to the codesopace, and fowarded to the Postgres server that running there and the results will be delivered back to my laptop.  But why is this useful?  We do want all our tools to be installed on the codespace, but I haven't found a postgres query tool that I like, so I want to run the Postico tool from my mac. I'll demonstrate this running in a later episode.

OK the next configuration value is `customizations`, and in there we have added some `vscode` specific `extenstions`.  VSCode extensions customise the VSCode editor and by adding them here those extensions will be available in ever codespace environment.  I've added the official `go` extension for VSCode which gives you syntax highlighting, code navigation, testing and debugging support. Its created by the go team and its really great.  The other extension is more of an experiment for me and adds support for editing "mermaid" diagrams inside markdown.

This is a good time to discuss documentation formats. I like markdown for documentation tasks because  it’s lightweight, easy to read in plain text, and automatically rendered with formatting on GitHub. If you look at the 'raw' versions of the markdown files you will notice that the diagrams are defined as text inside the markdown files.  These diagrams are defined textually using the mermaid language, which harks back to a point we talked about in an earlier video where we prefer to use textual formats for data, well these diagrams are just another kind of data. The beauty of this is that Github renders mermaid diagrams for you when viewing in the broiwser so its a super simple way that allows you to define diagrams in text, and also leverage Githubs ability to render them for you automatically. Also when someone edits a diagram you can see the differences in the textual representation.  

Right we are getting towards the end of this file.  We have two configurations entitled `onCreateCommand` and `postCreateCommand`, each of which refers to a separate shell script.

These scripts run at different times in the dev containers lifecycle:
- onCreateCommand runs **only once**, right after the container is **created for the first time**. It’s typically used for setup tasks that you don’t want to repeat, like installing dependencies, or downloading tools.
- postCreateCommand runs after **every container creation or rebuild**, once the container is up and running. 

In short: onCreateCommand = one-time setup; postCreateCommand = always after create/rebuild.

Lets take a look at what they run.

Open up the `on_create_command.sh`
The first line tells us its a bash script and it has some documentation at the top then it runs `make`

Which brings us to `make`, which is a Unix build automation utility that runs commands defined in a `Makefile` typically to compile applications or run repetative commands.

I've used make a Makefiles for many years, and the choice to use make os largely personal.  There are probably better tools out there, but at this stage I'm using what I know.  If you want to know more about make go to https://www.gnu.org/software/make/.  OK, so our command is `make setup` which means run the targer `setup` inside a `Makefile`.  Open up the Makefile and you will see some preamble at the top of the file, scroll down until you fnd the `setup:...` target.  When the command runs its going to find the target specified, run any dependencies, and then run the commands specified under the target. We have two targets that install tools via apt get and then install go tools.  Apt-get tools installed are the postgresql-client, git and jq. You can probably guess what postgresql-client is, it provides cli access to postgres databases, git allows us to run git commands for SCCS, and jq is a JSON query tool.  JSON is our data language format so we will use jq in a later video.  The go tool install command is a feature of go, which allows you to specify github projects as tools that you want to be installed, and when you run this command it will download those tools, compile them and install them in the codespace so they will be ready to run.  go tools come from the `tools` section in the `go.mod` file. Lets briefly go through each tool and what they do for us.

... discuss each tool ...

Ok, so at the end of the setup we have all the tools we need installed.

Remember all the tool setup happens once, as part of the 'onCreateCommand'. Now lets look at what happens in the 'postCreateCommand'.
 
Open the `post_create_command.sh` file.

The first non comment line calls set -e which tells teh script to exit immediately if any of the commands inside it return a non-zero status.

This may not be something you have encounytered before if you are new to a Unix based system. After running a command in unix, it will set a return code, which is an integer number.  The convention is that if the return value is zero, that means the command worked ok. Any other value indicates an error.

OK we have an if statement that checks for an environment variable called $GITHUB_ACTIONS, environment varaibles are dynamic values stored in the shell enviromnet and are used to communicate values across the system. The GITHUB_ACTIONS environment variable is set when we are employ the github infrastructure to automatically run our tests and we will cover it more in a later video.  

The -z option will be true if the GITHUB_ACTIONS envar is empty.  When the codespace starts up it will be empty so we will run the commands inside the if block. git pull origin will ensure that our codespace has all the latest code changes pulled down to the codespace environment.

OKm next we have a loop that iterates over any environment variable that is prefixed wity NATS_CREDS_. For each of those variables we base64 decode the value in it and save it to a file in our HOME folder. Base64 encoding and decoding is a standard way of sharing values as strings. The advantage of base64 encoding is that you can share any binary value, and it will be represented as an ASCII strings, which makes it perfrect for sharing in environment variables.  These values contain the secret credentials that we will use to connect to Synadia Cloud NATS.  You might ask yourself where those values come from.

Well in a later video we will sign up to Synadia Cloud, and create some credentials.  Then we will add them as codespace user secrets https://github.com/settings/codespaces.  Its super important never to commit credentials like API keys into our github repository's source files.  Luckily GitHub provide the codespace user secrets mechanism which allows us to save them in a safe place where no-one else can get them, and then have them automatically injected into our codespace at runtime.  So if users Jill, and Bob each create a codespace from the same github repo, they set up their own codespace user secrets and get their own values.

You will see inside the loop after the values are saved to a file, we call the `nats context add ...` comand. This registers the NATS credentils with the nats cli tool, so we can use them easily by 'selecting the context' later.

That marks the end of the startup sequence.

There is just one command left to run, which is `make bootstrap` which will run up all the services that need to be running for us to do our development.  In our case we need the postgres database to be running, and to have a two databases created, one for development called `beaker_development` and the other for unit tests called `beaker_test`.

Before we delve into the `bootstrap` target. I need to talk about docker because its the first time we have used it. What is docker?  Its an implementation of the Open Container Initiative standards, and it allows you to package up a full application including the libraries, code and other dependencies into a portable container format that can be run consistently across different environments like PCs, Macs, or server machines. Docker implements the OCI standard, but there are other implementations like PodMan and Kubernetes.  Our whole codespace environmnet is running inside an OCI runtime on Githubs infrastructure. We will be using Docker to run the Postgres database.  Docker is a key component in modern software architecture because its simplifies the way that we package software, for use on a wide multitide of systems. `docker-compose` is a tool that comes bundled as part of docker, it lets you define and run multi-container applications using a simple YAML file, so you can start everything (like a database, and email server) with a single command instead of running each container manually.  In our case we are only using it to run a database so it might be overkill. OK back to the `bootstrap` process

`bootstrap` depends on `restart-docker-compose`. The `restart-docker-compose` target, is dependent on two other targets `docker-compose-down` and `docker-compose-up` which it runs in that order. 

The `docker-compose-down` target will destroys and deletes our postgres server environment and all the data in its databases.  This approach of deleteing everything and starting from scratch is great for development and test environments where you don't have data that you need to retain. It forces you to understand your data requirements for each environment, and turn them into repeatable scripts. Once you have a scriptable environment, then its super easy to share with someone else, and have them configure or test their system exactly the same way as yours.   The `docker-compose -f .devcontainer/docker-compose.yml down  --remove-orphans || true;` command tells docker to use the `.devcontainer/docker-compose.yml` file and run the `down` command to stop all the containers defined in that file. The `remove-orphans` flag tells docker to kill any unconnected containers. We add the `|| true` on the end so the cleanup will allways run without error.  This is important because make will stop if any command returns an error, and at this stage we don't know the state of the docker environment so we just want to delete everything and take us back to an empty state. Ok so after this command runs lets just take it for a fact our docker environment isn't running any containers.  

Before we look at `docker-compose-up`, its time to examine the `./devcontainer/docker-compose.yaml` configuration file to see how we have scripted the definition of our postgres server. A YAML file is a human-friendly text file format used to represent structured data with indentation, often for configuration. This file lists the services we want to run inside docker.  We only have one service called 'db'. Inside that we have the 'image' which is the Docker image that we want run, the format is name:version, so you can see that we are running postgres version 18. The restart step tells docker to automatically restart the app if it crashes, but not if its explicitly stopped. The environment section contains a list of key/value pairs that represent environment variables passes on to the postgres image when it runs, and the port section is where we expose access to the application through TCP sockets.  TCP sockets work across machines over a network, using IP addresses and ports (like 127.0.0.1:5432), and are a bit slower than Unix domain sockets, but allow remote communication.  So when docker runs this image, its treated like a remote machine and any application running on the dev container will connect to it using a a TCP socket.  Each docker app has its own unique configuration settings, for example the Postgres app is documented https://hub.docker.com/_/postgres, which is a good place to look at the explaination for how its configured. This is where you would look to see what those environment variables mean.  I think they are pretty self explanatory, so if you want to know more I'll leave that as an exercise you can do online.

OK, back to the Makefile, after running `docker-compose-down` there are no docker containers running, so next it runs `docker-compose-up`. This creates all the containers specificed in our yaml file. the `-d` option means detach which tells the docker-compose command to run the containers in the background and return to us.  

Going back up to the `bootstrap` layer, now that we have recreated our postgres environment we run a couple of strange looking commands `$(MAKE) recreate-db DB_ENV=development` and `$(MAKE) recreate-db DB_ENV=test`.  All this means is make calls itself with a target `recreate-db` and an environment variable DB_ENV set.  So lets look at the `recreate-db` target. If you scroll down to find it in the Makefile it is dependent on 3 targets which it will run in order, `drop-db`, `create-db` and `migrate-db`, and then after they have run it uses the `echo` command to print a message to the terminal.  `drop-db` is dependent on `postgres-ready` and `terminate-conns`, after they have run it uses the psql command to `DROP DATABASE` using the DB_ENV environment variable to drop either `beaker_development` or `beaker_test`, but we are jumping the gun because we need to look at the dependecies first. `postgres-ready` uses the `pg_isready` command to wait for postgres to start-up, or wait 5 seconds if its not ready which will give it some time to become ready, then the `wait-for-it` tool waits for 30s for postgres to accept connections on port `5432`. This gives postgres some time to boosttrap and start its networking subsystem. At the end of this we know Postgres is up and ready to accept work in the form of SQL commands. `terminate-conns` uses the psql command line tool to run a SQL script that terminates any active connections to our beaker_test or beaker_development database. Why do we need to do that? Its because postgres prevents us from dropping a database that someone is connected to, so we force the connections to drop via this SQL. That makes life more convenient to the developer because otherwise you will have to find any processes that are connected to the database and shut them down manually.  Backing up a bit, now that postgres the database server is running, and we have created an empty database beaker_{DB_ENV} then the `create-db` target runs which will issue the SQL `CREATE DATABASE ..` command. This command creates a new, clean UTF-8 encoded database named beaker_<environment> (like beaker_dev, beaker_test, etc.), owned by the postgres user, with standard U.S. English locale settings. We will discuss some of these settings a bit more in a future video, but suffice to say we have created a fresh database with no tables or data in it.

That rounds out the `make bootstrap` command, but you might be wondering why we don't just run that automatically as part of the `post_create` startup sequence of the codespace.  I wanted to do that, but honestly I had a lot of trouble making it work reliably.  The problems seemed to stem from being unable to be sure that the docker daemon was fully started up. After spending quite a lot of time on getting this going, I just opted for a workaround which means the developer has to run `make bootstrap` when they start or restart the codespace.

This is something that happens often is real life development, you encounter a problem thats stopping you from achieving a bigger goal. Sometimes a manual workaround is fine if it helps us keep moving. We can always circle back later to tidy it up or automate once the pressure’s off.  The lesson here is that sometimes we can make a pragmatic call in the short term as long as we don't comprimise on the long term quality of what we are producing.  For me in this situation, my decision is that having an extra step for the developer to setup their environment is ok, because utlimately that wont comprimise the build of a production quality API server.


As our development environments become more complex, we need more tools and libraries installed, and we need more applications running in order to do development.  All this will slow down the creation of our codespace, which will impact our productivity.

But Github have a trick up their sleve that we can use called  prebuilt devcontainers.  A prebuilt devcontainer has the devcontainer image built in in advance to speed up container startup. Let me show you how thats done. Navigate to the github project, click on settings, the codespaces. You can see that I have a setup a prebuiult containere for us to use. If we click edit we can see how its set-up.  It runs only on the main branch and specifies the path to the devcontainer configuration to use. In the triggers section you can specify what will cause it to be triggered, and I've set it to be be when the devcontainer configuration changes.  In other situations, it might be better to do it periodically, say at 6am ever moning, so that when you team starts work they know there is always a fresh devcontainer ready to run. The other settings down below allow you to configure where the image will be available. I've set mine to australia, only so let me know if you are outside the area and are unable to run it. The last setting is a useful one, where you can specify to notify someone if your prebuild fails. 

OK lets click back and view the output to see how much time the pre-build step takes. Click on the ' See output' button shows that the pre-build took 30 mins. OK thats a signicant amount of time.  Lets now run up a new devcontainer in the browser and see what it looks like:

OK, so I'm going to run a new codespace in my Safari browser. I've also found that Chrome works well. Whatever browser your running I advise you to check its running the most up to date version before starting.  OK we start at the beaker homepage in github, Click code then codespaces, I'm going to click teh '+' biutton. It starts a new copy of VS Code, and in the terminal you can see its running the post_credate_command. After about 15-20s it tells us that its 'Finished configuring configuring the codespace'.  

Why does it start my codespace in VSCode and not a browser window? Thats because I can set a global preference against my github profile. Got to Settings > Codespaces > Editor prefernce to change that like I have.

Lets see what we have:

- On the left in VSCode the file explorer shows all the files that we have checked out of git.
- At the botton we cna see that we have checked out the 'main' branch, that can be changed, by the usual git commands, or can be set by creating the codespace when you have another branch selected in the browser window.
- Lets check if our tools are installed. Start a new terminal and lets check
  - `go version` confirms we have the go tools.
  - `psql -V` shows us the postgres client is installed
  - `jq -V` confirms that the jq is installed ok
  - Lets check some of the 'go' tools, we will just do a couple.
    - `sqlc version` is installed
    - `which wait-for-it` is installed

If something went wrong, you might wonder how we debug it. Weill is you were watching closely, you may have noticed a message ' Cmd/Ctrl + Shift + P -> View Creation Log to see full logs'. Lets so that now and have a look at the startup sequence. Scroll right to the bottom of the page and scroll up to see the last few steps.


The last step is manual, we need to run docker-compose script & finish creating a clean environment, by running `make bootstrap`
- Lets confiorm that ran OKby listing the running processes: `docker ps` shows postgres is running ok.

Ok we have covered a lot, so its time to summarise what we know:

- Codespaces runs in the cloud, and provides its user interface either through the browser or a local copy of VS Code. The minimum software requirements are a modern browser.
- We start a codespace on demand via the github web interface.
- There are two distinct phases our codespace goes through:
  - The 'on create' pahse runs exactly once. Thats our chance to install any tools we need. 
  - The 'post create' phase runs each time the codespace starts up. 


Right so that brings to to the end part of our video.  If you have opened a codespace then you can shut to down. To do that you close the window in the browser or in my case VSCode. The codespace is still running until you explicitly shut it down  or it shuts itself down after a period of inacivity.  The important thing to note is that if you have edited files inside the codespacem, they are retained until the codespace is fully deleted. So its perfectly normal to edit files one day, and save them, then restart your codespace the following day and pick up where you left off.

To delete a codespace, go back to the place where we created them. Click on the code button, then codespaces and the '...' button gives you the option to delete the codespace.

Github kindly offers a free quota of 120hours/month to use codespaces which is pretty amazing because it means we can all have a play with this amazing technology.

Lets summarise why we run our development environment using codespaces.
They give our teams consistent, reliable, stable environments for working in. Team members use them for tasks, discard them and create new ones quickly. The tooling is consistent across the team.  On-boarding new team members becomes a lot simpler because we can provide a reliable consistent development environment, by giving them access to Github and a browser.   The disadvantages is that it forces you down a particular set of tooling - most notably using VSCode. This might turn some developers off particularly is they love using a particular toolset thats not available with something like codespaces. 

We have learned how to create a new codespace, we know the lifecycle they go through, and how to use pre-builds to speed things up. When can debug a startup issue, and how we can keep things clean an tidy by deleting codespaces when we have finished with them

This brings us to the end of the video on 'development environments'   Thanks for listening , and remember "Iron sharpens iron, and one man sharpens another.”. Hit the subscribe button if you wnat to be notified when the next video is out. The next video in the series we start talking about data.

# Episode 4

