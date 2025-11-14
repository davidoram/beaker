# Transacript

# Espisode 1 - Introduction and overview

Hi and welcome to my series on "Production grade system development"

My name is Dave Oram, and in this series I'm going to walk through the elements of developing a production grade microservice API system.  My aim isn't just to show you code, I'll be discussing my decisions. Yes we will look at the code, but we will also cover important issues around design, developer environment setup, production grade runtime monitoring and system integration concerns.

For this reason the series is useful to CTOs, developers, architects, product managers, test engineers or anyone interested in the whole development process.

Each video looks at a different aspect of the system, and by the end of it we will have looked at every file in the project, so you should be able to get a sense for how the system works.  There will be lots of discussion and talk about what options were considered and what I chose and why.  These are just my opinions and what works for me. I've been developing systems for over 35 years, and what I'm presenting today models the systems I've set up in my day job. All these decisions involve personal choice and of course you should feel free to pick and choose the ideas and techniques that work for you.  The idea is to help develop some questions to ask yourself and your colleagues when making deicisons, so that you can feel confident in making good decisions for you.

OK, with that preambe out of the way lets get some basic housekeeping out of the way and get started.

All of the code that we will be walking through is publicaly available on github at https://github.com/davidoram/beaker. In a later video we will talk about the online development environment we will use, but for now you can just navigate to that page in the browser and take a look.  The project is MIT licensed.  This means you are free to borrow any or all of the code for your own use.  

OK, I've covered the introduction to this series, I look forward to seeing you in the next video where we cover off the 'High level architecture'.

Remember "Iron sharpens iron, and one man sharpens another.”. Hit the subscribe button if you wnat to be notified when the next video is out. See you next time.

# Episode 2

Hi and welcome to my series on "Production grade system development". My name is Dave Oram and I'll be your gude as we todays espisode which covers our "high level architectural goals".

If you want to learn more about this video series, I encourage you take a look at the previous video which provides a brief introduction to all the things that we will be covering.

OK, today we will be covering off the high level architecture. If you want to follow along point your browser at https://github.com/davidoram/beaker, or just grab a coffe, sit back and watch as we walk through the our high level architecture.

Starting at the README.md file which is displayed when you navigate to the github project.  There are a few statements in this file which I want to talk about.

The overview talks about a "production ready microservice API". What do I mean by that? Well there are a million tutorials on the web about creating APIs, but whats the difference between demo code and production code. I think of it as the difference between say a classic car, and a delivery truck. The classic car looks cool, but it needs a lot of TLC to keep it running well. Every few miles you need to top up the oil because its got a leak somewhere, you can only hold two people in the classic car because its a fancy convertable. On the other hand the truck is working day in day out. It needs to take heavy loads when necessary, and it needs to be super reliable.  There is a team of mechanics ready to jump in and fix the truck when it breaks because time is money with that expensive truck.  

 In the context of an API we might need to consider how that API fits into its environment. We want that API to be able to scale up when necessary and handle great load.  While that API is running we want tools to evaluate its peformance and health so we can be confident that its working correctly.  With a little bit of planning we can create a API that allows us to scale as our API gets more popular. In a future video we will be covering how we capture Telemetry information which provides real time diagnostic informtaion about our system as its running.

 In the next section I talk about three technology cornerstones used to build our system. The go programming language, the NATS messaging system and Postgres database.

 Why do I choose to program in go?  Well I've had many years experince with go so I'm well placed to compare it with other languages I've used in the past.  For me go is an excellent language for modern enterprise system development. There are many cool features built into go which you cna read about elsewhere, but I'm going to raise a few that affect productivity over the long term:
 
 - For starers, having a static typing eliminiates a whole class of run-time bugs, which for me means increased reliability. 
 - The tooling is second to none, and includes source code formatting, unit testing, linting, security vulnerability checking is built in.  This means teams can spend less time figuring out for example how to format code, or write unit tests, of check for security vulenabilities - those decisions are made for you so you can concentrate on writing code ro solve your business problems.
 - Go version compatability guarante. A new version of go is released every 6 months, you upgrade by recompiling and typically what that means is you get increased performance for free. Other languages evolve their syntax and libraries quickly, but it means you are constantly having to update code for no benefit apart from keeping current. And you need to keep your code current because production systems need to use the latest versions to stay ahead of security vulnerabilities. In my experience upgrading to the new version of go is a 10 minute task.

 OK lets click on the link to the "High Level architecture".

## High level Architecture

Lets look at all the pieces on the diagram in front of us.

OK we are building a microservice architecture.  What that means is we build small components, which helps you as a developer keep each component small and easy to understand.  We have boundaries between our microservices and if we keep the boundary interactions the same, then we can change the internal structure of any component without affecting the others.  As systems get larger that can be a great help in managing complexity. In this project we have a single micro-service for managing inventory.  If in the future we wanted to manage customers, we would make that a separate micro-service with its own application executable and its own separate postgres database. 

For the purposes of this exercise we will be running our microservice inside our codespace development environment. We will talk more about our development environment in a later video.

Our microservice stores its persistent data in the Postgres relational database system. Our database design and tooling will be covered in a separate video.

Our microservice code needs a way to communicate with the outside world, and we do that through NATS. NATS is an open source messaging system, and among its many capabilities it allows us to build API services which work on a request response model. To start with, all interactions with NATS are secured and require callers to be authenticated and authorized by NATS. To call an API a request is sent to NATS, that request is routed to the correct microservice which receives the request and it replys with a response that NATS routes back to the caller.   This is how regular HTTP APIs work as well, but in our case the transport layer isn't HTTP, instead our microservice uses the NATS protocol which is optimized for this kind of message exchange. What that means is that to call the API, by default the caller needs to use a NATS client rather than an HTTP client.  If you are an experienced developer at this point you may be asking yourself, hey HTTP is the lingua franca protocol for APIs on the web, so why would I want to use NATS and have my clients have to use a different protocol.

Well this is where Synadia come into the picture. They are the commercial company behind NATS and they offer a hosted NATS service caled Synadia Cloud. A hosted service means that you pay Synadia and they will run NATS on your behalf.  Lets take a look at what a Synadia Cloud gives us:

- The core NATS messaging platform, which includes Jetstream persistent messaging
- A Web UI Console for management
- HTTP Gateway
- Connectors
- Custom workloads

Back to the question of offering an HTTP API.  Synadia Cloud includes the 'HTTP Gateway'. Gateway is a common term that means a system or device that acts as a bridge beyween two different protocols.  In this case the HTTP Gateway converts incoming HTTP requests to NATS, calls the NATS service, and returns an HTTP response. Think of it as a 'translator' between the two protocols.

Using the HTTP Gateway doesn't require any extra work on your behalf, so its nice that you get dual protocol support for free.  If you are an enterprise customer, it's a killer feature because its often easier for customers to integrate with HTTP APIs rather than NATS because its so familiar to people.

Of course there are always tradeoffs that have to be made, and with the HTTP Gateway there are several to consider:

- You have no control of the URL, or the format of the API endpoint. Its always a  `PUT https://api.ngs.global/nats/subjects/{subject}` 
- You have no control, and no visibility over errors returned by the HTTP Gateway component. For example if your microservice goes offline then callers through the HTTP Gateway will get a `500 no responders` message.

OK, so like all the decisions we make across this project we weigh up the pros and cons and decide if this works for us.

For me, one of the main drawcards for using Synadia Cloud is that they provide, update and monitor all the infrastructure and settings needed to provide secure access to NATS.  At first glance it might seem easy, but here are just some of the things that they are doing for you when you use their service:

- They obtain and manage SSL certificates used by the TLS security layers to client connections secure.
- They provide a cluster of connection points to NATS, which is geographically distributed, and also multi-cloud. This provides redundancy and reliability
- They manage the Jetstream persistant message storage for us.
- Their teams are monitoring for security vulnerabilities, constantly patching the software stack, and also regularly upgrading NATS

Of course not everyone needs this kind of service, but if you are a small team, it can save your team a lot of time and effort if these things are already convered off. You can focus on building a great product.

You will notice that the inventory microservice publishes changes as events to Jetstream.  This is an example of an 'event-driven' architecture. In our case when something interesting happens inside the inventory service (Low stock event), so that is published as a event into NATS, and we have configured NATS to capture those messages in an persistent message stream.  Why is that useful?  It allows consumers to connect to that stream and read the messages.  Because the messages are stored to disc, if the consumer goes offline, it can re-join later and catch up on all the messages it missed while it was offline.  Think of it like a holding bay for those messages, where they can be stored temporarily until the consumers are ready for them.   This style of architecture where we offer actions orient APIs which allow callers to do something actively, in conjunction with events that can be passively consumed allows for you to build applications on top of these services in a very flexible manner.

The last piece of the diagram shows NewRelic as our OpenTelemetry provider.  OpenTelemetry provides a standardized way for your app to collect metrics, logs, and traces about its behavior and performance. That information makes it easier to understand, debug, and monitor applications. If multiple systems all use open telemetry their information can be connected together.  Then you can choose the vendor that suits you best to view all that collected information. There are lots to choose from https://opentelemetry.io/ecosystem/vendors/ so you can choose based on price, or performance.  I've opted to go with NewRelic because I use it in my day job, so I know how to navigate around its interface.

I think of OpenTelemetry like the dashboard of your car.  You have all sorts of guages and warning lights that help you understand how the engine (your application) is working. It will help us answer questions like, is the engine running, how fast id it workng, how much load is on the system, and is it about to break.

The last corerstone technoloy is the Postgres database. There have bene volumnes written about postgres, but needless to say its a high performance, very reliable RDBMS that we use to store our applications data. Its another open source system, that like NATS has many commerical offerings if we choose to rely on a third party to manage our database. We'll cover Postgres in more detail in a later eposode.

OK Thats enough about the high level architecture, lets focus in on the API that we are building:

## System Goals

If you scroll down on the same page to the "system Goals"  we will look at them one by one.

The goal is to build a basic inventory API that allows us to track quanity of different products.  It tells us that products are identified by a `product-sku`, and the three operation we have are to add, remove stock, and to query to see the current level.

The business rules section explains what constraints we want the system to maintain.

API endpoints explain each of the operations.

## Technical requirements 

Making the API accessaible via the NATS and HTTP protocol is straightforward. Weve already discussed why is useful to give two options, but I didn't explain why you might want to encourage your callers to try the NATS protocol option. With NATS you get the following advantages:

- Lower latency & higher throughput. the NATS client keeps a persistent connection and uses a lighter weight protocol.
- Bidirectional comms. Once connected the caller can send and recieve data from the NATS server. It natively support streaming style operations.
- Resiliance. The NATS client will automatically route requests across the cluster giving us better availability and fault tolerance.
- Security. With a NATS persistent connection the callers authenticate once and stay secuerly authentcated rather then having to authenticate each call.


Our API requests and responses will encode data using the JSON format. For those that don't know JSON is a human and machine readable format. https://en.wikipedia.org/wiki/JSON.  The Unix O/S environment demonstrates how we can build powerful applications by using text oriented tools.  Many people fall into the trap of worrying that their data isn't compact enough over the wire and they are losing a lot of performance by using a text oriented data format, but that concern should be balanced with how easy it is to start using your system, constructing input payloads and understaning responses.  Using a text format like JSON means we can eyeball the data and check very easily if things look right.  If we were to use a binary format like protobuf then we need special tools to help us do that, and each extra layer of tools slows us down.  The other advantage to textual formats is that other tools can manipulate the data from our system.  There is a large list of JSON tools and systems available https://github.com/burningtree/awesome-json

OK so JSON is our data format, but we need to layer some structure over that. We need to validate all out API inputs and very strictly control our API outputs.  To help us with that, we draw on another standard called JSON Schema https://json-schema.org. This allows us to define what our JSON data looks like and enforce validation rules.   This makes our validation and conformance to easy to implement.  With each incoming request we get the JSON Schema library to check if its valid before proceeding.  Like JSON, JSON Schema is an open standard so its used in tools like OpenAPI to define API specifications.  By adopting a standard we are able to integrate with other systems and tools more easily.

Next on the list we have a requirement to implement tests to verify our system is working correctly. This is always a good idea, because it gives us confidence that when we make changes they are working correctly and we haven't broken something inadvertatly.   We will go through our test code in a later video.

Our last requirement is to capture telemetry using the open telemetry standard. We have already discussed what open telemetry is, and we wnat to use it so we cna oberve our application becahious at runtime.

## Standards and Organisations

If you scroll down in the page you will see some of the projects that this application uses, either during development, or after its been deployed.  I added this so you can get a sense of what the projects are, and how they inter-relate.

Standards are useful because if they are widely adopted, and we make our applications follow them to then we get the benefits of joining a community.  For example I've used `git` and github as my source code control tool. Because gits a widely used tool, many organisations like github have standardised on its use, so I can take my git skills and use then inside my  github environment.

Lets look at the other tools, systems and standards.

I think its important to recognise that many many other people have freely shared their work, so I for one like to encourage and support them when I can. Anyone can do this by simply clicking the 'star' icon on the githib repository for their project.  Then spread the word with your colleagues and freinds so they can start using the awesome tools and systems that you find.

This briongs us to the end of our video on our "High Level Architecture", I look forward to seeing you in the next video where we cover off the 'Development environment'.

Remember "Iron sharpens iron, and one man sharpens another.”. Hit the subscribe button if you wnat to be notified when the next video is out. See you next time.


# Episode 3

Hi and welcome to my series on "Production grade system development". My name is Dave Oram and I'll be your gude as we todays espisode which covers our "development environment".

If you are new to this video series video series, I encourage you take go back and listen to previous videos as they cover some context to what we are covering today.

OK, today we will be covering off building a modern "development environment" in 2025. If you want to follow along point your browser at https://github.com/davidoram/beaker, otherwise you can just watch me cover all the steps.

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

Each cloud based development environment instance runs in the cloud, but looks and acts just like a normal development environment that you used to run on your laptop.

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
- Synadia Cloud - for our NATS service
- New Relic - for our Open Telementry data
- Github obviously for recording our source code changes.
- Not shown on the diagram, but codespaces has automatic port forwarding so for example you can run a website on your codespace and view it in your local browser. 


So what are the advantages of our coder environment:
- Your developer only needs a browser and an internet connection to get up and running. This makes onboarding a one click operation.
- IT managers might issue Chrombook $500, vs Macbook Airs ($1,800) vs MacBook Pros ($2,819). This offers significant hardware cost savings.
- Each developers environments is completely consistent on startup - because we script it
- Development environment setup is scripted so you can experiment and make changes with the knowledge you can back them out if they don't work.  You can not only change the software, but you can change the virtual hardware, increasing the number of CPUs or memory. 
- Enterprises can control access and costs
- Developers can run multiple independent codespaces concurrently. They autoamtically 'pause' and shutdown so you just pay for what you use.

But there are some disadvantages.
- Will developers accept a standardised set of tools. This is especially relenant to the IDE, so codespaces works best on VSCode and has beta support for JetBrains IDE, but what if your devs use VIM?
- It costs money. Personal accounts get a quota of free usage and its very afforable afterwards. For example I've been using codespaces for this project working in my spare time and my bill last month was 0.49 cents! Enterprise plans of course cost real money so that needs to be factored into the equation.


I've been using codespaces at my day job for over almost a year now and for me it works extremely well. You'll need to make your own judgements about how it works for you. Remember to “Seek wise counsel, but own the decision.”

OK, lets dive into a beaker projects codespace settings and examine  step by step how its put togther.

Click on the link to open the `dev-container.json` file. If you have checked out the project already it lives in thge `.devcontainer` directory in the root of the project.

Ok we start off with the `name` of our project which is informational only, so I called it 'Beaker project dev container'. There is nothing special about the name beaker - it was my nickname at a place I worked.

Next the `workspaceFolder` shows where all the application project files will live.

The `Features` section pulls in pre-packaged tools.  There are lots of them to choose from and I have added some links into the documentation.  This is the simplest way to start adding tools into our development environment. Ive added the `go` development environment and something called `docker-in-docker`.
The `go` feature adds the go programming language tools. `docker-in-docker` runs the docker daemon process inside docker, and allows our codespace environment to build and run docker containers.

The next section is `forwardPorts`. A forwarded port moves network traffic from one computer to another. In this context we are telling the codespace to forward network traffic from the codespace virtual machine to the developers laptop. Port 5432 is used to query the Postgres relational db.  So this means I can run queries from my laptop they will be re-directed to the codesopace, and fowarded to the Postgres server that running there and the results will be delivered back to my laptop.  But why is this useful?  We do want all our tools to be installed on the codespace, but I haven't found a postgres query tool that I like, so I want to run the Postico tool from my mac. I'll demonstrate this running in a later episode.

OK the next configuration value is `customizations`, and in there we have added some `vscode` specific `extenstions`.  VSCode extensions customise the VSCode editor and by adding them here those extensions will be available in ever codespace environment.  I've added the official `go` extension for VSCode which gives you syntax highlighting, code navigation, testing and debugging support. Its created by the go team and its really great.  The other extension is more of an experiment for me and adds support for editing "mermaid" diagrams inside markdown.

This is a good time to discuss documentation formats. I like markdown for documentation tasks because  it’s lightweight, easy to read in plain text, and automatically rendered with formatting on GitHub. If you look at the 'raw' versions of the markdown files you will notice that the diagrams are defined as text inside the markdown files.  These diagrams are defined textually using the mermaid language, which harks back to a point we talked about in an earlier video where we prefer to use textual formats for data, well these diagrams are just another kind of data. The beauty of this is that Github renders mermaid diagrams for you when viewing in the broiwser so its a super simple way that allows you to define diagrams in text, and also leverage Githubs ability to render them for you automatically. Also when someone edits a diagram you can see the differences in the textual representation.  

Choosing a simple text format for documentation reminds me that “Blessed are the simple, for they shall ship on time.”

Right we are getting towards the end of this file.  We have two configurations entitled `onCreateCommand` and `postCreateCommand`, each of which refers to a separate shell script.

These scripts run at different times in the dev containers lifecycle:
- onCreateCommand runs **only once**, right after the container is **created for the first time**. It’s typically used for setup tasks that you don’t want to repeat, like installing dependencies, or downloading tools.
- postCreateCommand runs after **every container creation or rebuild**, once the container is up and running. 

In short: onCreateCommand = one-time setup; postCreateCommand = always after create/rebuild.

Lets take a look at what they run.

Open up the `on_create_command.sh`
The first line tells us its a bash script and it has some documentation at the top then it runs `make`

Which brings us to `make`, which is a Unix build automation utility that runs commands defined in a `Makefile` typically to compile applications or run repetative commands.

I've used make a Makefiles for many years, and the choice to use make os largely personal.  There are probably better tools out there, but at this stage I'm using what I know.  If you want to know more about make go to https://www.gnu.org/software/make/.  OK, so our command is `make setup` which means run the targer `setup` inside a `Makefile`.  Open up the Makefile and you will see some preamble at the top of the file, scroll down until you fnd the `setup:...` target.  When the command runs its going to find the target specified, run any dependencies, and then run the commands specified under the target. We have two targets that install tools via `apt get` and then install go tools.  `apt get` is the unix installer tool. Apt-get tools installed are the postgresql-client, git and jq. You can probably guess what postgresql-client is, it provides cli access to postgres databases, git allows us to run git commands for SCCS, and jq is a JSON query tool.  JSON is our data language format so we will use jq in a later video.  The go tool install command is a feature of go, which allows you to specify github projects as tools that you want to be installed, and when you run this command it will download those tools, compile them and install the resulting binaries in the codespace so they will be ready to run.  go tools come from the `tools` section in the `go.mod` file. Lets briefly go through each tool and what they do for us.

... discuss each tool ...

Ok, so at the end of the setup we have all the tools we need installed.

Remember all the tool setup happens once, as part of the 'onCreateCommand'. Now lets look at what happens in the 'postCreateCommand'.
 
Open the `post_create_command.sh` file.

The first non comment line calls `set -e` which tells the script to exit immediately if any of the commands inside it return a non-zero status.

This may not be something you have encountered before if you are new to a Unix based system. After running a command in unix, it will set a return code, which is an integer number.  The convention is that if the return value is zero, that means the command worked ok. Any other value indicates an error.  

If you are new to software development on Unix its worthwhile studying its tools and the way they work because its good to “Remember the ancient paths, where the good way is, and walk in it.”. Even as an experienced developer I learned a lot from Eric S Raymonds book 'The Art of Unix Programming' it covers a lot of the thinking behbind Unix systems and its been a string influence on my thinking.

OK back to the code. We have an if statement that checks for an environment variable called $GITHUB_ACTIONS, environment varaibles are dynamic values stored in the shell enviromnet and are used to communicate values across the system. The GITHUB_ACTIONS environment variable is set when we are employ the github infrastructure to automatically run our tests and we will cover it more in a later video.  

The -z option will be true if the GITHUB_ACTIONS envar is empty.  When the codespace starts up it will be empty so we will run the commands inside the if block. git pull origin will ensure that our codespace has all the latest code changes pulled down to the codespace environment.

OK next we have a loop that iterates over any environment variable that is prefixed with NATS_CREDS_. For each of those variables we base64 decode the value in it and save it to a file in our HOME folder. Base64 encoding and decoding is a standard way of sharing data on the internet. The advantage of base64 encoding is that you can share any binary value, and it will be represented as an ASCII strings, which makes it perfrect for sharing in environment variables.  These values contain the secret credentials that we will use to connect to Synadia Cloud NATS.  You might ask yourself where those values come from.

Well in a later video we will sign up to Synadia Cloud, and create some credentials.  Then we will add them as codespace user secrets https://github.com/settings/codespaces.  Its super important **never** to commit credentials like API keys into our github repository's source files.  Luckily GitHub provide the codespace user secrets mechanism which allows us to save them in a safe place where no-one else can get them, and then have them automatically injected into our codespace at runtime.  So if users Jill, and Bob each create a codespace from the same github repo, they set up their own codespace user secrets and get their own values.

You will see inside the loop after the values are saved to a file, we call the `nats context add ...` comand. This registers the NATS credentials with the `nats` cli tool, so we can use them easily by 'selecting the context' later.

That marks the end of the startup sequence.

There is just one command left to run, which is `make bootstrap` which will run up all the services that need to be running for us to do our development.  In our case we need the postgres database to be running, and to have a two databases created, one for development called `beaker_development` and the other for unit tests called `beaker_test`.

Before we delve into the `bootstrap` target. I need to talk about docker because its the first time we have used it. What is docker?  Its an implementation of the Open Container Initiative standards, and it allows you to package up a full application including the libraries, code and other dependencies into a portable container format that can be run consistently across different environments like PCs, Macs, or server machines. Docker implements the OCI standard, but there are other implementations like PodMan and Kubernetes.  Our whole codespace environmnet is running inside an OCI runtime on Githubs infrastructure. We will be using Docker to run the Postgres database.  Docker is a key component in modern software architecture because its simplifies the way that we package software, for use on a wide multitide of systems. `docker-compose` is a tool that comes bundled as part of docker, it lets you define and run multi-container applications using a simple YAML file, so you can start everything (like a database, and email server) with a single command instead of running each container manually.  In our case we are only using it to run a database so it might be overkill. OK back to the `bootstrap` process

`bootstrap` depends on `restart-docker-compose`. The `restart-docker-compose` target, is dependent on two other targets `docker-compose-down` and `docker-compose-up` which it runs in that order. 

The `docker-compose-down` target will destroy and deletes our postgres server environment and all the data in its databases.  This approach of deleteing everything and starting from scratch is great for development and test environments where you don't have data that you need to retain. It forces you to understand your data requirements for each environment, and turn them into repeatable scripts. Once you have a scriptable environment, then its super easy to share with someone else, and have them configure or test their system exactly the same way as yours.   The `docker-compose -f .devcontainer/docker-compose.yml down  --remove-orphans || true;` command tells docker to use the `.devcontainer/docker-compose.yml` file and run the `down` command to stop all the containers defined in that file. The `remove-orphans` flag tells docker to kill any unconnected containers. We add the `|| true` on the end so the cleanup will allways run without error.  This is important because make will stop if any command returns an error, and at this stage we don't know the state of the docker environment so we just want to delete everything and take us back to an empty state. Ok so after this command runs lets just take it for a fact our docker environment isn't running any containers.  

Before we look at `docker-compose-up`, its time to examine the `./devcontainer/docker-compose.yaml` configuration file to see how we have scripted the definition of our postgres server. A YAML file is a human-friendly text file format used to represent structured data with indentation. YAML format is often used for configuration files. This file lists the services we want to run inside docker.  We only have one service called 'db'. Inside that we have the 'image' which is the Docker image that we want run, the format is name:version, so you can see that we are running postgres version 18. The restart step tells docker to automatically restart the app if it crashes, but not if its explicitly stopped. The environment section contains a list of key/value pairs that represent environment variables passed down to the postgres image when it runs, and the port section is where we expose access to the application through TCP sockets.  TCP sockets work across machines over a network, using IP addresses and ports (like 127.0.0.1:5432), and are a bit slower than Unix domain sockets, but allow remote communication.  So when docker runs this image, its treated like a remote machine and any application running on the dev container will connect to it using a a TCP socket.  Each docker app has its own unique configuration settings, for example the Postgres app is documented https://hub.docker.com/_/postgres, which is a good place to look at the explaination for how its configured. This is where you would look to see what those environment variables mean.  I think they are pretty self explanatory, so if you want to know more I'll leave that as an exercise you can do online.

OK, back to the Makefile, after running `docker-compose-down` there are no docker containers running, so next it runs `docker-compose-up`. This creates all the containers specificed in our yaml file. the `-d` option means detach which tells the docker-compose command to run the containers in the background and return to us.  

Going back up to the `bootstrap` layer, now that we have recreated our postgres environment we run a couple of strange looking commands `$(MAKE) recreate-db DB_ENV=development` and `$(MAKE) recreate-db DB_ENV=test`.  All this means is make calls itself with a target `recreate-db` and an environment variable DB_ENV set.  So lets look at the `recreate-db` target. If you scroll down to find it in the Makefile it is dependent on 3 targets which it will run in order, `drop-db`, `create-db` and `migrate-db`, and then after they have run it uses the `echo` command to print a message to the terminal.  `drop-db` is dependent on `postgres-ready` and `terminate-conns`, after they have run it uses the psql command to `DROP DATABASE` using the DB_ENV environment variable to drop either `beaker_development` or `beaker_test`, but we are jumping the gun because we need to look at the dependecies first. `postgres-ready` uses the `pg_isready` command to wait for postgres to start-up, or wait 5 seconds if its not ready which will give it some time to become ready, then the `wait-for-it` tool waits for 30s for postgres to accept connections on port `5432`. This gives postgres some time to boostrap and start its networking subsystem. At the end of this we know Postgres is up and ready to accept work in the form of SQL commands. `terminate-conns` uses the psql command line tool to run a SQL script that terminates any active connections to our beaker_test or beaker_development database. Why do we need to do that? Its because postgres prevents us from dropping a database that someone is connected to, so we force the connections to drop via this SQL. That makes life more convenient to the developer because otherwise you will have to find any processes that are connected to the database and shut them down manually.  Backing up a bit, now that postgres the database server is running, and we have created an empty database beaker_{DB_ENV} then the `create-db` target runs which will issue the SQL `CREATE DATABASE ..` command. This command creates a new, clean UTF-8 encoded database named beaker_<environment> (like beaker_dev, beaker_test, etc.), owned by the postgres user, with standard U.S. English locale settings. We will discuss some of these settings a bit more in a future video, but suffice to say we have created a fresh database with no tables or data in it.

That rounds out the `make bootstrap` command, but you might be wondering why we don't just run that automatically as part of the `post_create` startup sequence of the codespace.  I wanted to do that, but honestly I had a lot of trouble making it work reliably.  The problems seemed to stem from being unable to be sure that the docker daemon was fully started up. After spending quite a lot of time on getting this going, I just opted for a workaround which means the developer has to run `make bootstrap` when they start or restart the codespace.

This is something that happens often is real life development, you encounter a problem thats stopping you from achieving a bigger goal. Sometimes a manual workaround is fine if it helps us keep moving. We can always circle back later to tidy it up or automate once the pressure’s off.  The lesson here is that sometimes we can make a pragmatic call in the short term as long as we don't comprimise on the long term quality of what we are producing.  For me in this situation, my decision is that having an extra step for the developer to setup their environment is ok, because utlimately that wont comprimise the build of a production quality API server. Remember “God works through imperfect people — your code can too.”

So just to emphasis the point running `make bootstrap` is something we can do at any time to recreate our system with database that have their structure set but empty of data.  If we are experimenting and we break things its a quick and easy way to get back to a known state.  Environments like this encourage experimentation because we know that we can always revery things easily and quickly. "Not everything needs to last forever — just long enough to serve its purpose.”

As our development environments become more complex, we need more tools and libraries installed, and we need more applications running in order to do development.  Installing all these tools will slow down the creation of our codespace, which will impact our productivity.

But Github have a trick up their sleve that we can use called  prebuilt devcontainers.  A prebuilt devcontainer has the devcontainer image built in in advance to speed up container startup. Let me show you how thats done. Navigate to the github project, click on settings, the codespaces. You can see that I have a setup a prebuiult containere for us to use. If we click edit we can see how its set-up.  It runs only on the main branch and specifies the path to the devcontainer configuration to use. In the triggers section you can specify what will cause it to be triggered, and I've set it to be be when the devcontainer configuration changes.  In other situations, it might be better to do it periodically, say at 6am ever moning, so that when you team starts work they know there is always a fresh devcontainer ready to run. The other settings down below allow you to configure where the image will be available. I've set mine to australia, only so I hope that you can use it from where you are, if not maybe pop a comment on the video. The last setting is a useful one, where you can specify to notify someone if your prebuild fails. 

OK lets click back and view the output to see how much time the pre-build step takes. Click on the ' See output' button shows that the pre-build took 30 mins. OK thats a signicant amount of time.  Lets now run up a new devcontainer in the browser and see what it looks like:

OK, so I'm going to run a new codespace in my Safari browser. I've also found that Chrome works well. Whatever browser your running I advise you to check its running the most up to date version before starting.  OK we start at the beaker homepage in github, Click code then codespaces, I'm going to click the '+' biutton. It starts a new copy of VS Code, and in the terminal you can see its running the post_credate_command. After about 15-20s it tells us that its 'Finished configuring configuring the codespace'.  

Why does it start my codespace in VSCode and not a browser window? Thats because I can set a global preference against my github profile. Got to Settings > Codespaces > Editor prefernce to change that like I have.

Lets see what we have:

- On the left in VSCode the file explorer shows all the files that we have checked out of git.
- At the botton we can see that we have checked out the 'main' branch, that can be changed, by the usual git commands, or can be set by creating the codespace when you have another branch selected in the browser window.
- Lets check if our tools are installed. Start a new terminal and lets check
  - `go version` confirms we have the go tools.
  - `psql -V` shows us the postgres client is installed
  - `jq -V` confirms that the jq is installed ok
  - Lets check some of the 'go' tools, we will just do a couple.
    - `sqlc version` is installed
    - `which wait-for-it` is installed

If something went wrong, you might wonder how we debug it. Weill is you were watching closely, you may have noticed a message ' Cmd/Ctrl + Shift + P -> View Creation Log to see full logs'. Lets so that now and have a look at the startup sequence. Scroll right to the bottom of the page and scroll up to see the last few steps.


The last step is manual, we need to run docker-compose script & finish creating a clean environment, by running `make bootstrap`
- Lets confiorm that ran OK by listing the running processes: `docker ps` shows postgres is running ok.

Ok we have covered a lot, so its time to summarise what we know:

- Codespaces runs in the cloud, and provides its user interface either through the browser or a local copy of VS Code. The minimum software requirements are a modern browser.
- We start a codespace on demand via the github web interface.
- There are two distinct phases our codespace goes through:
  - The 'on create' pahse runs exactly once. Thats our chance to install any tools we need. 
  - The 'post create' phase runs each time the codespace starts up. 
- Our codespace makes the development environment consistent for everyone
- We can run on low cost hardware and tune our virtual hardware to suit our needs.
- I can shut down a codespace, open it up on a new machine and its exactly as I left it. It will retain all my uncomitted edits, open tabs etc.
- The disadvantage is that I can't work if GitHubs infrasturcure is down. That somtimes happens.

Right so that brings to to the end part of our video.  If you have opened a codespace then you can shut to down. To do that you close the window in the browser or in my case VSCode. The codespace is still running until you explicitly shut it down  or it shuts itself down after a period of inacivity.  The important thing to note is that if you have edited files inside the codespacem, they are retained until the codespace is fully deleted. So its perfectly normal to edit files one day, and save them, then restart your codespace the following day and pick up where you left off.

To delete a codespace, go back to the place where we created them. Click on the code button, then codespaces and the '...' button gives you the option to delete the codespace.

Github kindly offers a free quota of 120hours/month to use codespaces which is pretty amazing because it means we can all have a play with this amazing technology.

Lets summarise why we run our development environment using codespaces.
They give our teams consistent, reliable, stable environments for working in. Team members use them for tasks, discard them and create new ones quickly. The tooling is consistent across the team.  On-boarding new team members becomes a lot simpler because we can provide a reliable consistent development environment, by giving them access to Github and a browser.   The disadvantages is that it forces you down a particular set of tooling - most notably using VSCode. This might turn some developers off particularly is they love using a particular toolset thats not available with something like codespaces. 

We have learned how to create a new codespace, we know the lifecycle they go through, and how to use pre-builds to speed things up. When can debug a startup issue, and how we can keep things clean an tidy by deleting codespaces when we have finished with them

This brings us to the end of the video on 'development environments'   Thanks for listening , and remember "Iron sharpens iron, and one man sharpens another.”. Hit the subscribe button if you wnat to be notified when the next video is out. The next video in the series we start talking about data.

# Episode 4

Hi and welcome to my series on "Production grade system development". My name is Dave Oram and I'll be your gude as we todays espisode which covers our "Data layer and databases".

If you are new to this video series video series, I encourage you take go back and listen to previous videos as they cover some context to what we are covering today.

OK, today we will be covering off the "data layer and databases". If you want to follow along point your browser at https://github.com/davidoram/beaker, otherwise you can just watch me cover all the steps.

Its worth taking a take a lot of time and care to understand what our data is and how it inter-related, because thats the foundation of any appliction. Our data will live a lot longer than the applications that sit on top of them. Today we are building our application layer in go, but 15 or 20 years from now, we might be using another language, and we may need to bring our data from an old system into a new one. The code is re-written but the data is often migrated.


Data is at the core of any application. In our scenario, we are building a stock keeping API, so our data is in the form of Product SKU's and quantities.

The process that I start with when modelling real world data is to map them onto standards and domains.

I don't know much about Product SKUs so I searched online to see if there was a standard for representing SKUs. I couldn't find one so I looked around and designed a data type that represented what I thought was reasonable, a string between 1 and 64 chars long, containing alphanumeric underscore and hyphens. My thinking here is that spaces are probably not generally used, so we want don't allow them. Why? because if I was to represent two different SKUs "a1" and "a1 " on screen it would be very hard to tell they are different.

To take another example, imagine if our system modelled telephone numbers we have standard E.164 - and its defined https://en.wikipedia.org/wiki/E.164. If I wanted to store a phone number I would follow this standard. Why?  Because I know that if the data from system needs to interact with another system, its much more likely that they will talk together. Data interchange between systems is an important measure of the usefulness of a system, so adopting standards will help your system interact with other systems more easily which is added value. There are standards for all kinds of data type, for example names, addresses etc. I encourage you to find some online - wikipedia might be a good place to start looking. Many applications appeal to a global audience, so if you data structures can cope with that right from the start you are at a real advantage.

The other data attribute we will be modelling is stock levels.  After consulting with the business experts for my system they let me know that there is no upper bound for inventory levels, but they can't ever fall below zero. We also don't sell fractions of a SKU, so stock levels are best modelled as integer values that are >= 0.

Because we started by gaining a good understanding of our domain data, we can now use that information to map that data into our application.  

Our lowest level of data management is Postgres so start with that. Why did I choose Postgres over some other tool.  The main reason is that I have many years experience with it, its safe reliable, and available everywhere. There are many other valid alternatives but Postgres is a good choice for me because I can explain the features that I'm using as we go through the code.  The bible says “Store up for yourselves treasures in heaven, where moth and rust do not destroy.”, and Postgres fits those same values for me because its very safe and reliable. 

In a previous video we showed parts of the Makefile where we start a Postgres server using docker compose.  I choose the latest version at the time of writing which is 18. 

If you open the Makefile and find the `create-db` target it shows how we create the database. Lets go through the options and what they mean.

- WITH OWNER postgres: Makes the PostgreSQL role postgres the owner of the new database.
-	ENCODING 'UTF8':	•	Ensures the database uses UTF-8 character encoding.
- LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8': 	Defines locale settings for string sorting (collation) and character classification (ctype). This means:
  - Sorting and comparisons are case-insensitive for ordering purposes in the sense that "A" and "a" are considered equal in sort order (they’ll group together).
  - Equality checks are still case-sensitive in standard SQL. 'a' != 'A'
- TEMPLATE template0. Creates the new database as a copy of template0, which is a “minimal” empty template database in Postgres.	This ensures no extra extensions or locale settings are inherited accidentally from template1.

OK so we have a an empty database how do we create the tables, indexes and our data definitions inside the database so its ready to be used. 

We use a tool called `sql-migrate` to apply database migration files which are basically versioned changes to the database.   So how does it work, we write migration in plain SQL files, then the sql-migrate tool applies them in the correct order.  It tracks which migrations have been run in a special database table. The end result is scripted database setup that we can maintain just like any other source file. `sql-migrate` runs against databases other than Postgres, and provides other tools to reverse out migrations etc.

OK, lets look at running the migration. We actually did it as part of the last video but we can do it again by running `make bootstrap` which will recreate the database, and run the migrations. When we run that `sql-migrate` reports that its `Applied 1 migration` Lets look at the sql-migrate command the first option `up` means apply any migrations that haven't yet been applied. The `--config dbconfig.yml` option points to the configuration file for sql-migrate and the third option --env specifies an environment which in our case will be either `development` or `test`. Lets take a look at the config file. In there there are two sections, one for each environment and under that options that tell sql-migrate the `dialect` which dictates database specific behaviour, `datasource` which tells it how to connect to the db, `dir` for migration files, and `table` that contains the record of which migrations have run.

Lets look at our `db-migrations` directory to see what files it executed.  There is only one `db-migrations/20250716085349-create-tables.sql`. The strange filename uses the date+time smashed together, and becomes a convenient ordering mechanism thats used to order the files, and when you have more than one file sql-migrate will apply them in the same order that vs code shows them.  Open that file and you see some SQL commands and some comments.  First of all the `+migrate Up` and `+migrate Down` comments tell sql-migrate which block of SQL to run.  We are migrating `up` so lets focus on the top block.

Our system has a single database table called `inventory` with two columns. The `product_sku` has type text which is the appropriate data type to hold strings, and we have made it not null, and the primary key. By making it not nullable we prevent any NULL values from being stored in that column, so we can only have strings, and by making it the primary key Postgres will automatially add a unique index onto that column, enabling super fast lookup when looking for an exact match. 

Next the `stock_level` column is defined as an integer also not null.

Now we have database constraints. A database constraint is a rule that the database enforces on a table or column to ensure data integrity — basically, it stops invalid or inconsistent data from being entered.

Our inventory_stock_level_nonnegative ensures stock_level is never negative. If someone tries to insert `-5`, the DB will reject it.
inventory_product_sku_format ensures product_sku only contains lowercase letters, numbers, hyphens, or underscores, and is 1–64 characters long. Prevents invalid SKUs like `ABC!@#` from being stored.

Constraints like these are your first line of defense for keeping the data correct and predictable. With constraints, the database itself guarantees correctness no matter how the data gets inserted or updated, by the application or a migration script. The rules are defined in one place and are often very efficient because they are run inside the database server right next to the data itself.


OK lets check the database is working OK.  I'm going to use a tool called Postico which is a third party app that can connect right from my mac into the Postgres database.  I can do that because in a previous video we saw the devcontainer was configured to `forwardPorts` 5432 which makes that possible to connect from my mac -> devcontainer -> postgres running in docker

Lets test the inserts - see 'test inserts'


Lets load 100k rows into our database & check how query performance works

```sql
DO $$
DECLARE
    batch_start INT;
    batch_end   INT;
BEGIN
    FOR batch_start IN 1..100000 BY 1000 LOOP
        batch_end := batch_start + 999;

        INSERT INTO inventory (product_sku, stock_level)
        SELECT
            'sku_' || gs::text AS product_sku,
            (random() * 99 + 1)::int AS stock_level  -- random int 1–100
        FROM generate_series(batch_start, batch_end) gs;
    END LOOP;
END $$;
```



Select a row `select * from inventory where product_sku = 'sku_58991';` Takes ~50ms which is perfectly fine for development.  If we had a "real" production grade Postgers server running say in AWS or Google cloud this would be much faster.

OK that concludes the lowest layer of our data management which is Postgres, now we are going to go up a layer and talk about how we manage that data in our application.

Our application is written in `go`, so the first thing we need to do is connect our go application with the postgres database. 

We are going to start walking through the code or our application but I'm only going to focus on the parts relating to database connections. Lets start in `cmd/main.go`.  

In the imports you will see `"github.com/jackc/pgx/v5/pgxpool"`. So pgx is a database driver that allows go apps to run SQL against Postgres. Its not the only one but its one that I've used a lot. Its very mature and reliable which is why we are using it.  Note  that we are pulling in the `pgxpool` module which provides a connection pool.  

When designing a high performance application, it might get many requests at once. Each request needs a database connection, so we want to control or limit the amount of concurrent connections, so we don't overload the database by using too many connections.  pgxpool is a  connection pool for the pgx PostgreSQL driver in Go. It improves application performance by maintaining and reusing a pool of open database connections instead of opening and closing a new one for every operation. The pool automatically creates and destroys connections as needed, and allows us to limit the maximum number of concurrent connections in use at any one time.

The first line `ctx, cancel := context.WithCancel(context.Background())` is somewhat related because this context permiates througout the code.  In Go, a context is used to manage deadlines, cancellation signals, and contain request-scoped values across API boundaries and goroutines. It allows you to control the lifecycle of operations—such as shutting down gracefully, timing out, or propagating cancellation—especially in concurrent or networked applications. Contexts help coordinate work and resource cleanup, making your programs more robust and responsive to external events. We will touch on contexts as we walk through the code.

The next line is `setTimeZoneToUTCOrExit` which sets the timezone of our application to run in UTC stands for Coordinated Universal Time. It’s the primary time standard the world uses to keep clocks and time consistent. Think of it as the “baseline” time zone with no daylight savings, no regional offsets — just a stable reference clock. The rule of thumb when developing apps is to process and store any times in UTC and convert to local time only when presenting to the user. Now our app doesn't store any times but its an important consideration for most apps that store when this happen.

Clicking through the GetOptions takes us to where the command line arguments are parsed.  We always have the database connection details passed in at runtime when the application starts, so our app can connect to different databases, but in the same way. We have one option which is the --postgres which is a URL containing the connection string that encodes the host and port to connect to, username and password along with the database name and any other options.  Using this format allows us to connect with low security to our local dev or test database, but then switch to better security or different options when connecting to a production server. See https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING

Back to the main code and the next block of code is `setupPostgresPoolOrExit` and it passes in a context, the call to ParseConfig allows the connection pool to be configured straight from the URL, but then I override this to have a maximum / min connection limit and a connection lifetime.  This means that connections once created will be discarded after they have been used for 30 mins. We then create the new Pool with that configuration and "ping" the database to check we are connecting ok. Finally we return the pool.

Back to main, we call defer `pool.Close()` so that the application cleans up when it exists. 

The pool is then passed in when we create the new `App` struture.

Now that we have established our connection pool,  we are going to skip a whole bunch of code and focus in on how we use the databae inside the app.

For our purposes, each API call request needs to interact with the database, so its going to take a connection from the connection pool, start a database transaction, run some queries and then if everything works ok commit the transaction and return the connection back to the pool for use by the next request.  If things go wrong instead of committing the transaction, we will perform a roll back operation which will undo any changes made to the database.

This is a common pattern which helps us achieve the goal of Atomicity, whereby our API call might perform 10 queries agianst the database, some inserts, updates and deletes but they should happen atomically so we wrap all of the operatings in a database transaction and ensure either all of them happen together or none.

Now we are going to examine the code we use to actually interact with the database so its time to look at the `sqlc` tool.

sqlc is a developer tool that generates type-safe code from SQL queries.
- You write plain SQL files (e.g., SELECT * FROM inventory WHERE product_sku = $1;).
- sqlc reads those queries and your database schema, then generates Go functions and structs.
- The generated code handles parameters and results safely, so you don’t have to hand-write query boilerplate. The typesafety aspect ensures that the correct mapping between go types, and Postgres types is maintained.

The benefit is sqlc lets you keep full control of your SQL while getting the convenience and safety of strongly typed code in Go. Tools like sqlc are worth considering because boilerplate or generated code saves you time, is completely consistent and more reliable than hand-crafted code.

To use sqlc it needs access to the database to read the table definitions, and we get some configuration options when it comes to code generation.  All that is done through the  `sqlc.yaml` file, so lets open that and take a look.

- version: "2": The config file format version — 2 is the current stable format.
- sql:The list of SQL generation targets. Each item describes how to generate code for a given schema and set of queries.
- schema: Path (or directory) containing your schema files — CREATE TABLE, ALTER TABLE, etc. sqlc uses this to understand your database structure and types.
- queries: File or directory where your SQL query files live (e.g., SELECT, INSERT, etc.). sqlc reads these queries and generates matching Go functions.
- engine: Specifies which SQL dialect to use — here it’s postgresql.
- gen: → go: Tells sqlc to generate Go code.
- package: – the Go package name for the generated code (db).
- out: – output directory where generated code will be placed (internal/db).
- sql_package: – which Go SQL driver to use (pgx/v5 instead of the standard database/sql).
- database: → uri: A connection string to your local or dev database.
sqlc connects to it to validate queries and infer correct types.

OK, so now we have our configuration sorted out lets examine our queries.  Open `query.sql` and examine the queries we have. All are prefixed by a specially formatted comment that sqlc will use to determine what kind of function it needs to generate. the `name:` part will end up being the go function name, and the suffix `:one` tells sqlc what the return value will be.  In all of our queries we are returning a single row, but this could be `:many` to return multiple rows. You can also specify `:exec` if no resultset is expected. There are many option options here see https://docs.sqlc.dev for more details.

Lets examine each query in turn. First our AddInventory query. Add this product to inventory, and if it already exists, increase its stock instead of creating a duplicate. This is an example of a single query that performs an upsert. It trys the insert a new row with the product_sku and stock_level. The on conflict is triggered if there already exists a row with that product_sku, in which case instead of failing the query it will update the existing stock level. In this part of the query `EXCLUDED.stock_level` refers to the value being inserted and `inventory.stock_level` refers to the current value in the table.

This kind of query helps improve system performance. We know that inserting and updating stock is a common operation, so we want it to happen as quickly as possible. By performing one SQL operation instead of two we halve the number round trips between our application server and our database.  When we examine telemetry in a later episode we will be able to quantify exactly how long these operations take.

The other two operations RemoveInventory and GetInventory are straigtforward.

One other thing to note is the use of $1, $2 which represent the parameters pased into the query.

Once we have setup the configuration file and written our queries, then we run the sqlc compiler, to generate the go source code that our application will use. To do that we can run `make sqlc` which calls `sqlc generate` to generate the go code. Lets look at the files in `internal/db/` to see what has been created.

Lets start with `db.go` It provides a `Queries` interface that our application will use, and a `DBTX` interface to be used by the generated code that will use to execute SQL. Note that the `Exec`, `Query` and `QueryRow` functions are tailored specifically to the `pgx` database connection library, if we had chosen the `lib/pq`, or `database/sql` standard library as our `sql_package` this would look diferent.  To create a `Queries` struct we can call `New` and then optionally call `WithTx` to wrap multiple sql calls in a databse transaction.

There are two other files created. First lets look at `models.go` sqlc has generated structs that represent the tables we are querying. In our case there is only one the `Inventory` struct that contains fields for the each column. Note how `sqlc` has automatically mapped the columns to the correct data types in go.   If you define your own data types, or have third party data types you can configure sqlc to use those data types. 

The other file created is `query.sql.go` open that up and lets look at whats been generated..

The first thing you notice is that each query we defined has been turned into constants. Where a query takes more than one parameter sqlc defines those paramaters in a struct as we see with `AddInventoryParams` which has defined the two paramaters using their correct go types.  

Lastly we see the function that our app will use. They all follow the same form, you need a `Queries` object and they take a standard context, and query arguments, and they return the structure representing the row returned and an error. The `AddInventory` function calls the  `QueryRow` and then using the already defined the Inventory variable, casts the response values into it and returns.

The other functions work similarly.

OK so that wraps up the 'data layer and database' discussion   Thanks for listening , and remember "Iron sharpens iron, and one man sharpens another.”. 

Hit the subscribe button if you wnat to be notified when the next video is out. The next video in the series we conitinue our discussion about data, but from the perspective of how our API consumes data and validates it

# Episode 5

Hi and welcome to my series on "Production grade system development". My name is Dave Oram and I'll be your gude as we todays espisode which covers our "Data and APIs".

If you are new to this video series video series, I encourage you take go back and listen to previous videos as they cover some context to what we are covering today.

OK, today we will be covering off the "Data and APIs". If you want to follow along point your browser at https://github.com/davidoram/beaker, otherwise you can just watch me cover all the steps.

In the previous episode we discused the importance of modelling data in our database, and how we use the sqlc tool to get typesafe access to that data from the appilication.  Todays video covers how data is represented in our API layer, how to model it and validate it.

Our API is one of the most important parts of our system. Its the interface as seen by the outside world, and a key part of our API is its representation of the underlying data.

In an earlier video we briefly discussed that our API will use JSON to represent data. JSON is used in three key ways:
- For API requests - the data sent by the caller
- In API responses - the data returned by the API back to the caller
- When publsihing data  in events. for example we publish an event when inventory is low

So we have estalished that JSON is a great data format, because :
1.	Human-readable – easy to read and debug.
2.	Lightweight – small, simple text-based structure.
3.	Language-independent – supported natively in almost every programming language.
4.	Structured – represents nested data (objects, arrays) clearly.
5.	Easy to parse – most frameworks have built-in JSON parsers and serializers.
6.	Widely used – it’s the de facto standard for REST and modern web APIs.

OK great, but how can we impose our domain specific JSON structures over JSON?  Thats where JSON Schema comes in. We will use it to validate JSON data, Document our expected requests and responses. In short: JSON Schema ensures our JSON data is well-formed, consistent, and predictable across systems.  It's good to remember the timeless wisdom “Test everything; hold fast to what is good.”

We will use JSON Schema files to define our requests, responses and events.

We are not going to quite as far as generating an OpenAPI definition for our API, but thats more certainly possible because we are using the same building blocks as used by that standard. The team at my day job does just that.

What tools will we need to use JSON and JSON Schema?  Well go has good support for JSON, so thats an easy tick. It has a package that allows us to marshall and unmarshall JSON to and from go structs which we will use.  For the JSON Schema part we will use the third party tool https://github.com/santhosh-tekuri/jsonschema.  This is both a library and a cli tool, and we will be using both.   

The cli tool will be used in our build phase, to check that our JSON Schema files are well formed. 

At runtime we will verify that the request coming in via an API matches the JSON Schema before processing it, so in effect it performs validation on the request. This saves us a lot of development time, which otherwise we would have had to do manually.  Because JSON Schema allows you to define validations inside the schema we can use the library to perform those validations for us. We effectively move the validation from go code -> JSON Schema configuration files. Writing less code is allways a good goal to strive for so this fits with our philosophy.

I vaildate that API responses match the schema in unit tests.

Lets start by looking at a JSON Schema file. Open up the request for stock add in schemas/stock-add.request.json.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
```
- Declares the JSON Schema version (draft 2020-12). This is the latest version of the JSON Schema standard.

```json
  "$id": "http://github.com/davidoram/beaker/schemas/stock-add.request.json",
```
- Unique identifier (URI) for this schema so it can be referenced by others.

```json
  "title": "stock-add.request",
```
- Human-readable title for the schema.

```json
  "type": "object",
```
- The JSON data must be an object (key-value pairs).

```json
  "properties": {
```
- Lists the valid fields inside the object.

```json
    "product-sku": {
      "$ref": "http://github.com/davidoram/beaker/schemas/product-sku.json"
    },
```
- The `product-sku` field must follow another schema that defines valid SKU values.

```json
    "quantity": {
      "type": "integer",
      "minimum": 1,
      "description": "The number of units to add, must be at least 1."
    }
```
- The `quantity` field must be an integer ≥ 1, with a helpful description.

```json
  "required": ["product-sku", "quantity"],
```
- Both fields are mandatory.

```json
  "additionalProperties": false
}
```
- No extra fields beyond these two are allowed.

**Summary:**  
This schema defines a “stock add” request object that must include a valid `product-sku` and a positive integer `quantity`, and forbids any other properties.

But we dont yet  fully understand what the `product-sku` definition is. Its defined by `$ref` `http://github.com/davidoram/beaker/schemas/product-sku.json`, which tells the JSON schema compiler to look elsewhere for the definition of product-sku.   

Lets open up `schemas/product-sku.json` to look at how that is defined.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
```
- Declares that this file uses the 2020-12 version of the JSON Schema standard.

```json
  "$id": "http://github.com/davidoram/beaker/schemas/product-sku.json",
```
- Gives this schema a unique identifier (URI) so other schemas can reference it.

```json
  "title": "Product SKU",
```
- Human-readable name for the schema.

```json
  "type": "string",
```
- The value must be a JSON string (not an object or number).

```json
  "pattern": "^[A-Za-z0-9_-]+$",
```
- A regular expression that defines valid characters — only letters, digits, underscores, and hyphens are allowed.

```json
  "description": "The SKU (Stock Keeping Unit) identifier for the product. Consisting of alphanumeric characters, underscores, or hyphens. Min length of 1 character. Max length of 64 characters.",
```
- Explains what this field represents and its constraints, useful for documentation tools.

```json
  "minLength": 1,
  "maxLength": 64,
```
- The string must be at least 1 character and at most 64 characters long.

```json
  "examples": [
    "SKU12345",
    "PROD-001",
    "ITEM_67890",
    "a",
    "A1B2C3D4E5F6G7H8I9J0K1L2M3N4O5P6Q7R8S9T0U1V2W3X4Y5Z6_7-8"
  ]
```
- Example valid SKU strings that conform to the schema.

**Summary:**  
This schema defines a product SKU as a string between 1 and 64 characters, containing only alphanumeric characters, underscores, or hyphens — ensuring SKUs are consistent and valid across the system.

OK, so thats great but how does the JSON Schema validator tool link the two files? There is a big clue in the `schema-lint` Makefile target.  When you run `make schema-lint` it iterates over each *.json file in the `schemas` directory, passing each to the `jv` command line tool.
The `--assert-format` option tells jv to check that the file being validated is correctly formatted as a valid JSON Schema, not just valid JSON syntax.
The `--draft 2020` option specifies the JSON Schema version (draft 2020-12) to validate against.
Passing `--map http://github.com/davidoram/beaker/schemas=schemas`
Creates a URI mapping — so when the schema references something like
`"$ref": "http://github.com/davidoram/beaker/schemas/product-sku.json"`,
`jv` knows to resolve that reference locally from the `schemas` directory instead of fetching it over the network.

OK so when you run `make schema-lint` it should say "All schemas are valid", but lets double check that by making sure it picks up an error.  Lets edit the `schemas/stock-add.request.json` and introduce a typo into the `product-sku` $ref. When you re-run `make schema-lint` it will fail eith an error.

So how do we perform checks inside the API at runtime.  It involves the following steps.
- Creating JSON Schema files
- Setting up a JSON Compiler
- Performing the validation of the raw JSON
- Defining go structs to match the schemas
- If validation works then marshallinbg the raw JSON into the go structs so the application now has that data available to it.

Lets start going through the steps.

We have already seen the JSON Schema file definitions, for example `schemas/stock-add.request.json`. We factor out any common types into their own files so we can re-use them - like we did for `schemas/product-sku.json`

Next lets look at how the JSON Schema compiler is constructed. Open the `internal/utility/schema_loader.go` and we are going to walk through the `NewJSONSchemaCompiler` method.  It takes a context (which is unused - my bad), and a schemaDir, the folder that holds the schema files.  We pass that to the `NewLoader` function passing a map from the url prefixes to the directory that holds the schemas with that prefix.  The loaders job is to resolve schemas as its validating. When the schame compiler needs a specific schema as identified by its id, it asks the Loader to retrieve it.  Our Loader will load our applications schemas from that directory, but then it will also the HTTP loader to retrieve schemas from across the internet.
It creates an HTTPLoader which is a wrapper around the standard library http.Client setting the timeout to 15s. Then we create a FilePrefixLoader which checks if the url prefix matches one thats mapped to a folder, and if so loads that file from the disc & returns it. If no match is found the request is passed on to a jsonschema.SchemeURLLoader which delegates to other loaders based on the url scheme - thats the starting portion of the url eg: `file:` or `http:` etc, so if its a file it will use a FileLoader that simply loads files from the local file system, if its an http/https scheme it delegates to the HTTPLoader which issues an HTTP GET request to retrieve the file & return it.

Now a little sidebar here. You might be thinking hey I don't see you caching any of these requests. Thats right. I always avoid caching until is proving to be a problem. Why? Because caching is hard. Caching is hard because you’re trading correctness for speed — and keeping that trade fair is tricky. You need to know when the truth changes (invalidating the cache), You need to know how long to cache for (TTL), in a distributed like the web there are many layers of caches - which layer is lying to us. All this makes for a tricky debugging scenario. I allways like to verify if what I have is "fast enough" before I start looking at caching. "Do the simple things first; God handles the impossible refactors.”

OK, back to the `NewJSONSchemaCompiler` function, now that the loader is created, we create a NewCompiler, and attach the loader. Then we call `AssertContent` which means if the schema says this field should contain a specific type of encoded content (like base64 data or a certain media type), double-check that it really does, and `AssertFormat` which means if the schema says this field must be an email address, make sure it really looks like one — not just any string. Some JSON Schema compilers treat these things as hints but don't enforce them, and we want to be as strict as possibile with our inputs. Lastly we tell the compiler to use the latest draft version of JSON Schema. 

When we want to use the compiler we will call  the `ValidateJSON` function in the internal/api/request_scope.go file. The first few lines capture telemetry data, which is important to know later - we will look at that in a later episode. Incidently capturing precicisely how long an operation like schema validation takes allows us to make data driven informed decisions. If the idea of caching schemas comes up you can look at the telemetry and guage exactly how loing is spent doing this work and decide if you need to optimize that part of the code path or not. Don't be fooled into thinking faster is better everywhere, sometimes you trade-off so much complexity to gain a millisecond or two it may not be worth it when you can look elsewhere for easier optimizations.

OK, we check if rs.HasError - lets come back to that in future. We next check id the jsonData is emapty - thats immediately an error.  If the compiler hasn't been passed in thats also an error.  Finally we get to ask the compiler to `Compile` a schema for us based on the schema name we have asked for.  Again this is another potential place for caching, but I don't need to so I keep the code path simple. If the schema is nil, thats an error, but if not we use the function that the library provides to unmarshall the JSON, again checking for errors. Now we get to the important part. The call to `schema.Validate(data)` - this takes the JSON, and validates it against the schema returning an error if it doesn't comply. Its a simple as that one call to perform all the multitude of validations that your JSON Schemas can apply.

It means I can change my schema files to include a new field, I specify that field is an email I don't have to change one line of my validation code because its all encoded in the JSON Schema file.  I encourage you to take a look at the web. Some examples incliude:
- Regex: https://json-schema.org/learn/miscellaneous-examples#regular-expression-pattern
- Complex objects with nested properties: https://json-schema.org/learn/miscellaneous-examples#complex-object-with-nested-properties

The key point I'm making is that we move validation from code into configuration files. So whats the trade-off. Theres always a tradeoff. What is means is that the error handling might not look as nice as a custom error message because the validator generates those error messages.

Right lets assume our incoming data matches the JSON Schema, how do we work with in inside the `go` applictaion code.  Thats actually pretty simple, we marshal the data from JSON into a go struct using the standard library. 

Lets take the example of our incoming 'stock add' request. The go struct definied in file `schemas/stock_add_request.go` and its super simple. Open up that file and you will see near the top we define a constant representing the schema `StockAddRequestSchema`. We will be using that later in the application code when we ask the jv library to validate some JSON against a specific schema.  Next it defines the `StockAddRequest` struct which contains the same fields we defined in the JSON schema file. Whats important to note is the use of structure tags. Those are the notataions inside the backticks to the right of each field, and they tell the standard librarys json module how to map JSON into the struct (called unmarshaling), and from a struct back out to JSON (called marshaling). You can learn more about struct tags in go here https://go.dev/wiki/Well-known-struct-tags. Its something that your code can use to attach metadata to struct fields that can be extracted and used at runtime.

Once we have marshaled the data into a struct its able to be used in the application code.  

Once the application has finished working with the data , we need to trasform the go struct back to JSON and return it to the caller. 

Once we populate a go structure representing a response like `schemas/stock_add_response.go` with the data it needs,  we marshal it from go struct -> JSON and return the response. This is done through the standard library function json.Marshal, which uses the same struct tags used by the unmarshaling process, but this time it converts the data in the struct tag to JSON representation.

The return structure shows an interesting feature in the go struct tags. You will notice that the `ProductSKU,Quantity,Error` are all pointers and the struct tags might say "ProductSKU *string `json:"product-sku,omitempty"`"  The `omitempty` directive will tell the marshall process to skip that field if its empty or nil. It will become more obvious when we look at the guts of the API handlers in the next video how we can utilise this to make our life easier.

So lets take stock of what we have learned. 

- JSON is the chosen format for API requests, responses, and events due to its readability, lightweight nature, and wide language support.
- JSON Schema is used to validate, document, and enforce the structure of JSON data for requests, responses, and events.
- Go's built-in JSON support is leveraged for marshalling/unmarshalling, while the jsonschema library is used for schema validation both at build time and runtime.
- Validation logic is moved from Go code to JSON Schema files, reducing manual code and improving maintainability.

OK so that wraps up the 'data layer and API' discussion   Thanks for listening , and remember "Iron sharpens iron, and one man sharpens another.”. 

Hit the subscribe button if you want to be notified when the next video is out. The next video in the series we will discuss Telemetry. I'm looking forward to seeing you then.


# Episode 6

Hi and welcome to my series on "Production grade system development". My name is Dave Oram and I'll be your gude as we todays espisode which covers our "Telemetry".

If you are new to this video series video series, I encourage you take go back and listen to previous videos as they cover some context to what we are covering today.

OK, If you want to follow along point your browser at https://github.com/davidoram/beaker, otherwise you can just watch me cover all the code.

I've touched on Telemetry before but in this episode we are going to do more of a deep dive  into what it is and how essential it it to wring production grade systems.

Telemetry is a fancy word that means 'measuring something' and sending it somewhere for monitoring.  Think of in a hospital the heart monitor attached to a patent is sending a stream of data to the nursing station.

Its a bit like that with our application. We want to be sure that our system is online and responding correctly, so that we can serve our customers well.  Some systems define these requirements more formally in something called 'service level agreements' that might state the response times for an API, eg: 95% of responses in 100ms and 99% in 1s.  So we can use telemetry to measure that capture that particular metric and grade our system against its requirement.  Like the nurse monitoring the patient we need that information to be available and measurable to evaluate their health.  Another thing that a heart monitor might do is trigger an alarm when a the heart rate goes outside of the normall range. The alarm might be an audible alert designed to draw attention quickly to an urgent problem.  We can use the same technique when our systems are behaving outside their expected bounds.

We are using the standard called Open Telemetry that allows us to write our code to a vendor independent standard, and then choose from a number of vendors based on price, usability and other factors.

Our architecture here is that we will have out application send telemetry directly to our vendor of choice. This is a bit of a cheat on my behalf to keep things simple. In a real production envrionment we would more likely run a piece of software called a collector or agent. Our app would send its telemetry to the collector, which in turn can be configured to filter, aggregate, and apply custom rules, before it sends telemetry data to the vendor. 

Right so lets talk about the vendor I've chosen. We are going to sign up for New Relic using the instructions provided [here](./otel.md#setup-with-newrelic).  NewRelic, and many other vendors kindly offer a free account which is very generous, and allows you to perform some testing 


If you are new to this video seBefore we run it lets take a look at that target.  It executes the `otel-cli` tool, which is a third party tool that provides the abaility to integrate with an OTEL back end. You will notice that at the top of the Makefile we setup an OTEL_ENV which configures all the standard environment variables that OTEL tools use. This includes the endpoint of the NewRelic back end service, the API key, the service name, and attributes to pass with each trace. The `otel-cli` tool when you use the `exec` option measures the time it takes to run the command. It also notes the commands return value, and if its non-zero it will record that as an error on the span. Run `make test-otel` a and also run `make test-otel-error` so we can simulate some traces.  The major difference between these targets, is the command it runs.  The `test-otel` runs the unix sleep command to sleep for a second, and the `` target runs the `false` command which simulates an error because it always returns non-zero.  

Note you might have to give NewRelic a few seconds to ingest process and display your traces.  This is normal.  On the `traces` page if defaults to showing us traces from the last 30 mins, so that will include our traces, but you can adjust thatest-otel-errort and go back in time. It shows the number of traces, spans, duration and errors which get little graphs. We get some good infomation on this front page showing errors, average trace duration, etc.

Click on the `test-otel` name and you get a list of traces - we only have 1, each with the name, when it occured, the number of entities, spans and errors. Lets click on one and view some details. The one yellow bar shows our trace which has a single span. The trace is a single logical operation, and spans are sub-steps.  Later on we will see some example but for now our operation has a span with 1 trace.  Its 1s long, and when I click on it I can see more information pop up in a panel to the right. It starts with some "performance", which aren't useful for our test case but might be when we have thousands and we are looking at one. It might help us see some trend - eg: this span might take significantly longer than the average. Click on the next tab called "attribnutes" this is a key one.  When we create traces and spans we can attact attrinutes which is user defined data to the span to capture extra meaning. OTEl published [semantic conventions[(https://opentelemetry.io/docs/specs/semconv/) which has a [go library](https://pkg.go.dev/go.opentelemetry.io/otel@v1.38.0/semconv/v1.37.0) with hundreds of methods to help you follow the correct naming conventions. For example if you want to follow the convention for setting your "service name" you would use the https://pkg.go.dev/go.opentelemetry.io/otel@v1.38.0/semconv/v1.37.0#ServiceName function. Like other parts of the system, following the conventions established by the OTEL standard means that our system will work with a wide variety of tools. If we use the standard names, then we increase the chance of our system working with standard tools.

We are now going to run our service and examine a trace.  I know we haven't looked at that code yet but we will get to that in the next video.  For now its important for us to see what a real API call/response looks like so we have a sense for what to look out for.


We are going to run  Makefile targets that use the `nats` cli tool to exercise the endpoints. They authenticate as a stabdard caller, send a request and cpature the response, formatting that response nicely through the `jq` tool so the output in color so its easier for us to read.  

Lets run `make run` to start our API server in one terminal. In a new terminal run `make test-get` to call the `stock.get` endpoint.

When the response comes back we can see it includes some text lines before the JSON response.  It tells us the subject where the request is sent, the rtt (round trip time), which in my case is ~300~ms which is pretty slow, and finally the `traceparent` header value.   

Lets open NewRelic and find the trace using that **traceparent** value

This trace looks a little different from the test traces we examined earlier. First of all the 'entity map' diagram shows our 'beaker' app and 'postgres'. A newcomer looking at this knows that this trace interaced through the 'beaker' application which interacts with a database.  Clicking the expand all button displays all the spans that make up this trace. remember a span is a smaller pice of work inside a trace. Lets look at them individually:
- 'setup db conn' wraps two child spans that acquire a database connection and begin a transaction. Together this takes 0.76ms - pretty fast
- 'validate JSON' Performs the JSON Schema validation and it takes 0.4ms - really fast
- 'decode request' which converts the JSON string to go structs is even faster at 0.02ms
- 'add stock' wraps a query, that INSERTs into inventory.  The query itself has two child spans. Clicking on the query shows the SQL to the right. Executing this query takes around 1.4ms
- Committing the query takes 0.9ms.
- Building the response is very fast at < 0.01ms
- Sending the response takes 0.05ms

OK, What does this trace tell us?
- The server side time to process this is approx 4ms - which is fast
- There are no errors.
- About 50% of the time is spent in database operations. 

But from the clients perspective it took 319ms - so whats going on in that balance of about 315ms thats not taken up in out API.

- The client is establishing a TCP connection to the Synadia NATS global service endpoint
- Synadia NATS has to authenticate and authorize the API call
- The message is send
- The response is decoded and displayed.

Remember that our client and server are a very low powered codespace environment. In a production environment we can get better performance by:
- Running our API server in the clous (say AWS) rather than in a codespace
- Running a separate OTEL collector
- Using a high performance Postgres server hosted in the cloud say AWS
- Using a long running Client connection, rather than reconnecting each API call, so the cost of AuthN is amortized across all calls.

I would expect to be getting < 10ms API reponses for something like this in a production setup.

OK, great so now we have some basic information about our API calls.

Lets talk about when things go wrong.  Back to the analogy, where the nurse looking after the pateint. By regularly looking at traces we can see how our system is performing over time.  But now we want immediate action when some important event occurs, ie: when the alarm goes off.


In our context a few critical events can occur:
- One scenario is when we our system has an error, which we simulated earlier with our `test-otel-error` Makefile. 
  - In this case we definately want NewRelic to alert us when any traces come through that have been marked as having an "error". Systems like NewRelic have integrations that allow you and your team to receieve alerts, through email, slack or teams. Or maybe even open tickets for you make them visivle with your other work.
- The other might be a requirement from the business, it might be of great interest when stock levels fall below a threshold, because we need to re-stock those items.  
  - One option is that you could handle those things by recording metrics in our Telemetry system. Thats a different part of Open Telemetry that I'm not going to go into with this series, but there is plently of information about that online.
  - Another option is that our application can deliver its own feed of messages when that happens. In fact thats exactly what we do, when the  `stock-remove` endpoint is called and stock levels fall below 10. It will publish a message to a NATS subject. I'm pointing this out to let you know that telemetry isn't the only way for our system to notify adjacent systems of changes, in fact publishing messages is a powerfull technique for building loosely connected apps that react in real time to changes. But thats a topic for another time.

OK so that wraps up our 'Telemetry discussion. Thanks for listening , and remember "Iron sharpens iron, and one man sharpens another.”. 

Hit the subscribe button if you want to be notified when the next video is out. The next video in the series we will finally getting into the guts of how our microservice pulls all these threads together and implements our API handlers. I'm looking forward to seeing you then.

# Episode 7

Hi and welcome to my series on "Production grade system development". My name is Dave Oram and I'll be your gude as we todays espisode which covers our "Microservice implementation".

If you want to follow along point your browser at https://github.com/davidoram/beaker, otherwise you can just watch me cover all the code.

In previous episodes we have covered our database layer, how we use Postgres to store data and sqlc to write boilerplate access code, we've covered off the adoption of JSON as a data interchange format, and JSON Schema as the schema definition language for describing our API requests and responses. JSON Schema handles our data validation. Go supports unmarshall JSON text -> go structures so our app can use it and then later marshall from go structs back to JSON text again so that data can be returned by the system. Finally we touched on Telemetry and how we use that to record system behaviour.

In todays video, we will draw those threads togther and show how the NATS service framework can be used to build microservices.  

The NATS open source system has great docs to help you start building services at https://docs.nats.io/using-nats/developer/services

Lets cover what we mean by a service:
- A service has a group of logically related functions
- Services are discoverable. Meaning there is a way to query the system and discover what services are available
- Services have one or more endpoints that, which represent operations that the service provides.

OK, so lets talk specifics for our service.
- We have one service called "beaker"
- Indside "beaker" we have three endpoints:
  - "stock-add"
  - "stock-remove"
  - "stock-get"


Before being allowed to call any API we require callers to prove who they are and that they have permission to call our API. We call this Authentication and Authorization.  In our case these functions are delegated entirely to NATS. NATS has this functionality built in, its well designed by security experts so we can be confident that its a solid foundation to build upon.

**Sidebar** Using NATS to solve our AuthN/AuthZ unburdens us from having to do implement these functions.  This is very hard to implement right, and it feels like a good decision in many circumstances because we are letting experts implement this function. This isn't just my recommendation. Take a look at Microsoft Secure Coding Guidelines, or the OWASP recommendations or RFC 7435. Its generally accepted that its poor practice to roll your own security.  However we must consider the downside.   Although NATS will provide AuthN/AuthZ when a request is routed to our service we have no idea who made the call. We just have to trust that NATS has checked they are allowed to do that.

OK, so back to the NATS Services - lets talk about how they work.  NATS is a messaging system, and as such it supports a request/reply messaging pattern. This coventiently matches exactly what our microservice wants to do, the caller issues a request, the service decodes and processes that request and responds with a reply.

How is the request routed to the correct microservice endpoint?  This is where we define a unique NATS **Subject** for each service to listen for requests. Subjects are strings that form unique names or addresses that publishers and subscribers can use to find each other.  

Remember earlier I mentioned that services must be 'discoverable'. So there are some well know Subjects that all NATS services use to share information about their endpoints and the nats tools know about them so they can use that to find out the list of services and their endpoints. Thats explained at https://docs.nats.io/using-nats/developer/services#service-operations 

Now that we understand how services get their requests, how do they get their reply? When the request is sent the caller behind the scenes generates a unique Subject for the reply to be sent back ok, that only that caller will know about and be listening on.  It will look something like `_INBOX.lcWgjX2WgJLxqepU0K9pNf.mpBW9tHK` where the `_INBOX` prefix is static, but the rest is random and unique.  So when the caller receives a request, its given the Subject on which to send the reply.  All of this detail is hidden when write our code using NATS Service framework, but its immportant to know because our first step is to create a NATS user with all the permissions needed to act as our Microservice.  We mentioned this an earlier video when we setup our codespace, and now its time to set this up.

If you recall we are going to use Synadia Cloud as our NATS service provider. Our API caller will connect to the hosted NATS as will our microservice, and NATS will route the incoming requests and outgoing responses between them.

So lets head over to Synadia and sign-up for their free plan. Navigate to https://www.synadia.com/cloud and click on the "Get started for free" button. Just make sure you are signing up for "Synadia Cloud" because they have a few product offerings. I signed up through GitHub which allows me to sign-in through that which is super convenient. Once you are in there you will be presented with a list of **Systems** which has only 1 called NGS (NATS Global System), click on that and it shows a list of **Accounts**. 

Each Account is like its own separate namespace or environment. Users live **within** Accounts so they can only connect to that Account. This makes Accounts partiularly useful as a tool for SAAS account separation. By default NATS Users in one Account can't communicate with another Account.  

OK, for our test we are going to create some users in the "default" Account, so click on that Account then click Users. 

We are going to create three NATS users and make them all available for use in out codespace.

- The first is called "App" and its the user that our microservice uses to connect to nats. Create that user with default permissions which gives it the rights to publish/subscribe to **any** subject. In a real production scenario you would adopt the principle of least permission which means you only give this user the least permissions possible for it to do its job.
- The next caller is called "Caller" and its the user that we will give to our end user to call our microservice. We **will** adopt the principle of least permission and only give them. Add them with the following permissions:
  - pub: `$SRV.>`, `stock.>`
  - sub: `$INBOX.>`

- The last caller is called "CLI" and its our user that we will use to do anything within the API from the command line.

For each of these callers:
- Click get connected 
- Download the credentials file
- Transfer to codespace
- Type in in terminal `base64 /path/to/credentials-file`
- Copy the result to clipboard and go to https://github.com/settings/codespaces. Then add then each as secrets called `NATS_CREDS_APP`, `NATS_CREDS_CALLER` and `NATS_CREDS_CLI`
- Delete the credentials file.

**Important** don't forget to delete the credentials files from your codespace, so you don't accidently commit them to your git repo.  If you did that accidently you can just 'revoke' the credentials from the synadia UI & create some new ones.

OK, so lets test them out by starting a **new** codespace, which will allow the codespace to pick up the new variables.

When the codespace is up, we can `View > Creation Log` and scroll to the bottom and see that our codespace startup sequence has automatically added the NATS credentials to our file system, by decodeing the values in those environment varibales, and turning them into files. Then it called `nats context add ...` to save them into the nats context so we can try them out at any time.

Lets start by setting the using `nats context select` and selecting `NATS_CREDS_CLI`. This user can basically do anything.

Type `nats service ls` to list the services running, and there are none. On Synadia cloud website lets view the "Connections" tab to see which Users are connected. There are none because the `nats ...` command,, connected, issues its command and disconnected.   So now I'm going to start out microservice running by staring a new terminal and calling `make bootstrap run` (We run bootstrap because we allways need to do that when we start a open a new codespace).  You might see some error here around not being able to send telemery data - just ignore those for now,

 It takes a minute or two the first time because it has to download libraries etc. But once thats done it should say something like `"INFO beaker is running"`. When that happens go back to the Synadia cloud web UI and you will see 1 connection 'beaker'. Clicking on it shows useful metrics on the right. It shows the IP address, and connection details. Below that is shows the NATS Account & User. Lastly it has subscriptions and statistics.

 Before we dive into the code, lets illustrate how we make an API call.  Open the `Makefile` and find the `test-add` target. It uses the `nats` CLI tool to connect using the context we set up earlier called `NATS_CRED_CALLER`, which determines the authentication and authorisation rules.  Next it submits a `req` which is short for issueing a standard request/reply, the next parameter `stock.add` specified the NATS Subject where the request is sent, and the final paramater is the JSON payload which specifies the product-sku and quantity. I skipped the `--translate=` option which passes the result through the `jq` command line tool to format and colour the output so its easier for us to view.

 The results show something like:

 ```
06:37:59 Sending request on "stock.add"
06:38:00 Received with rtt 321.262441ms
{
  "ok": true,
  "product-sku": "coffee-cup",
  "quantity": 10
}
```

Ok, so lets review how this call is made:
- The request starts off in our codespace, initiated by the `nats` cli tool.
- It Authenticates against `tls://connect.ngs.global` which is Synadias endpoint for the hosted global NATS supercluister
- The User is Authentiocated and Authorized against my 'Default' NATS Account, and the request is allowed to be sent to the NATS subject `stock.add`. The payload is just a bunch of bytes, and the sender has also supplied a random response Subject eg: `_INBOX..lcWgjX2WgJLxqepU0K9pNf.mpBW9tHK`. The nats cli is listening for the response on that Subject.
- The beaker server we have running is listening for requests to the `stock.add` Subject on that same NATS Account, so it received it, decodes the JSON, and processes it by interacting with the Postgres database, and encodes the response into JSON and sends it to the `_INBOX.lcWgjX2WgJLxqepU0K9pNf.mpBW9tHK` Subject. The response is just a bunch of bytes.
- The `nats` cli receives the response and translates the output through `jq --color-output .` which colorizes and formats the JSON response.

OK. So we know know at a high level whats happening. Lets look at how we implement one of these services.  We will focus on the `stock.add` endpoint that we just called.

In a previous video we walked through `cmd/main.go` and looked at how the application starts up. It sets the Timezone to UTC, it configures our telemetry, creates a pool of Postgres connections, connects to NATS, creates a JSON Schema compiler. It sets up a signal handler then it starts the app.  On all of these bootstrap functions if something goes wrong the application exists with a non zero exit code.  This is the unix standard way of saying the an application failed for some reason.  Anyway if all goes well lets see what happens in the `startAppOrExit` func, which just calls `api.StartNewApp()`

The `App` state is represented a struct that has links to all the subsystems that the app needs, a NATS connection, Database pool and JSON schema compiler. It also holds a micro.Service which is the NATS standard definition of a microservice.   Lets navigate to the `makeService` function to see how that works.

It first creates a `micro.Config` struct that defines the service.  This information will be available when a caller introspects all the services available in a NATS system. The Service has a Name, Version, and Description. We also add an ErrorHandler which will be called in the case of an unhandled error, which just logs the error.

Calling `micro.AddService(...)` registers the service.
Next we add a `Group` called "stock" which simply groups a bunch of endpoints having the Subject prefix 'stock'.  To that Group we add the `stock.AddEndpoint(...)` handler. The first argument is the suffix to add to the group's subject so "add" becomes "stock.add" Subject.  The next argument is the Handler function that will process requests. The actual handler is implemeted in the `stockAddHandler` func, but I wrap it in a `traceHandler`.  Lets look at the `traceHandler` first - its further down in the same file.  

The traceHandler function wraps a microservice request handler, starting a new OpenTelemetry trace span for each request, logging the API request with context, and then invoking the original handler with the traced context. This enables distributed tracing and contextual logging for each API endpoint call. We discussed telemetry in the previous video, so for now lets focus on the fact that it simply logs each incoming request, then calls the handler.

Open the [internal/api/add_stock.go](internal/api/add_stock.go) file to see the `stockAddHandler` function.  It has the App struct as its method reciever, so it can access the database, JSON Schema compiler etc.

At the  top of the file we define `stockAddScope` and we need a small sidebar to discuss this.

```go
type stockAddScope struct {
  *requestScope[schemas.StockAddRequest]
}
```

This is a tiny helper type that wraps the generic `requestScope` for the specific `StockAdd` request. By using this wrapper we get a simple, concrete receiver for our handler methods. Those methods can call `rs.Request()` to access the decoded `schemas.StockAddRequest` value. The pattern keeps handler code easy to read while still reusing the same generic scope implementation.


It works as follows:
- First it creates a NewRequestScope, passing in the context, request, NATS conn, and db pool
- Next it calls defer rs.Close to ensure that when this function returns that request scope will be properly cleaned up
- Then the request scope  is directed to Validate the incoming request passing in the context, json schema compiler, request data (JSON payload) and a JSON Schema name.
- The request is Decoded into a `schemas.StockAddRequest` struct
- Then AddStock is called with that request, and the response from that used to MakeStockAddResponse
- Then the changes are  CommitOrRollback
- Finally the RespondJSON function is called to send the response.

So it's an 8 line function that starts by creating a `requestScope` and then calls a handful of helper methods that operate on that scope.

So we should start by describing what a `requestScope` is. Let's examine the request-scope.go file.

The `requestScope` struct holds all the state relating to a single API request. It has:

- a connection with nats
- the incoming request
- an error
- A connection from the postgres connection pool
- a database transaction
- a db.Queries object.
- the decoded request value of type `T` after calling `rs.decodeRequest(ctx)`, for example our `StockAddRequest`


OK, so when the `newRequestScope` func is called it creates & returns a new requestScope. It calls `setupDbConn` passing the connection pool. Lets take a look at that func. It performs some telemetry - we will skip over that for now. Next it calls `pool.Acquire` which gets a free connection from the connection pool. This is the first thing that can fail with an error so lets look at how errors are handled.  

The overall design philosopy is that the request struct holds the first error that occurs inside its struct.  When **each** step of the request processing starts it begins by checking - is there an error on the request & if there is it skips its normal processing.

OK, so when an error occurs we can call one of two functions against the request struct. Either `AddSystemError` if its some error that occurs inside the **system** that the caller has no control over.  "system" errors are the developers responsibility to monitor for, and fix if needed.  The other kind of error are called "**Caller**" errors and they are the responsibility of the API caller to fix. There is a function `addCallerError` to add them.

So in our case if the error is caused when we acquire a connection from the db connection pool, thats a **system** error so we call that function.  In turn that calls the `addEror` func that performs some telemetry, creating spans and logs which we will talk about later. Then it simply stores the error on the `requestScope`.

Back up to the `setupDbConn` lets assume the happy path and we get a connection, it saves the connection in the `requestScope` and proceeds to creates a database Transaction, perform similar error handling and stores that also against the `requestScope`.  Lastly it creates a new `db.Queries` and stores it against the request scope.

So back up to the `newRequestScope` function, we now have populated the request itself from nats. This gives us access to the incoming request data. We have a nats connection  saved (so that later on we can send a response), and we have established a connection to our database through our `queries` field. The queries are using a connection gained from the pool, and wrapped in a db transaction so we can commit or rollback as required later. 

Right lets return to the `stockAddHandler` function.
The first thing we do is defer a call to `requestscope.close` Defer of course means that it will be called when the function returns, so lets return to the close function until we have looked at the rest of the this function.

The first thing we do is call the `validateJSON` func passing in the context, JSON Schema compiler, the request data (the JSON payload in []byte form), and the name of the schema that the payload shoud confiorm to (ie: `"http://github.com/davidoram/beaker/schemas/stock-add.request.json"`).

validateJSON, creates a span for telemetry purposes.  Then it follows our standard check which says if the request scope has an error, return. Asumming all is ok, we perform some checks.
- First verify the the jsonData isn't empty - if so here we add a Caller error - one thats entirely caused by the callers actions. and return.
- If the JSON Schema compiler isn;t initialised thats a system error. Agian we record that and return.
- Now we ask the compiler to compile the schema name, which returns an schema object we can use to validate the incoming data. If that fails its a system error and we return
- Finally we unmarshall the bytes to an object using the `jsonschema.UnmarshalJSON` function. If its invalid JSON this will return an error.
- Finally we ask the `schema` object to `Validate` that data, which confirms that the JSON matches the Schema definition.  This is where our library performs all the "heavy lifting" validating all the elements of that JSON for us. If it fails at this step its a "caller" error and we bail, but if it suceeds we have good JSON input and the `validateJSON` method returns.

OK, back to `stockAddHandler`, next we call `decodeRequest`. This function looks a little different from the others because its a **generic** function. I'm passing in `schemas.StockAddRequest` inside square brackets, and that parameter tells the function what type I want returned by the function, it has regular parameters of a context and the request scope object. Lets look at how it works.

The`decodeRequest` function is a generic helper that decodes incoming request data into a specified type T, while integrating tracing and error handling:

- It starts an OpenTelemetry trace span named "decode request" for observability.
- If the requestScope already has an error (checked via rs.hasError()), it immediately returns. This follows our standard pattern of not continuing if something failed earlier
- It declares a variable decodedRequest of type T, which will hold the decoded result.
- It attempts to unmarshal the raw request data (rs.req.Data()) from JSON into decodedRequest.
If unmarshaling fails, it records the error as a "caller error" in the requestScope (using rs.addCallerError).

Recall in an earlier episode that the `json.Unmarshal` populates our go struct using struct tags.  Now we have our input request as native go type sitting on our request scope, and we are in a position for the rest of the request handling code to access it.

The next line calls `rs.addStock(ctx)` first, so lets look at that function.This function works as follows:

- It starts an OpenTelemetry trace span named "add stock" for observability.
- Then we get the request struct and use it to create the `db.AddInventoryParams`, setting the product sku and quantity. If you recall from our earlier episode on the database layer, `sqlc` has generated this structure for us. 
- Then we call the `rs.queries.AddInventory()` func that will actually call the SQL query to insert / update the inventory level 
- If the query fails, we add a system error & return nil
- Lastly assuming all goes well, it returns the new inventory level, a structire of type `db.Inventory` again generated by `sqlc`.

Back to the `stockAddHandler` func, we pass the result from `addStock` into the `makeStockAddResponse` function whos job it is to generate a response. Lets see how that works

- It starts an OpenTelemetry trace span named "build stock-add response" for observability.
- Then it creates a new response struct
- Next it makes a decision. If any error has been recorded up to now we say the response is ok = false, and we attach the error string
- Or if no errors detected up to now we set the `ok = true` and build up the response product SKU and Quantity
- Lastly we just return the response.

Back to the `stockAddHandler` func, we call `commitOrRollback`. This function works as follows.

- First it does a defensive check - if the request never had a transaction created we cant have ever created a transaction, so we just return straight away.
- Then we create a msg which is either `tx commit` or `tx rollback` depending on if an error was recorded, and we starts an OpenTelemetry trace span with that msg for observability.
- Then we create a function that will set the request scope transaction to nil when the function returns.
- Next if the request scope reported an error we issue a rollback, otherwisw we issue a commit. In both cases we add a SystemError if that fails.

Lets return now to the `close` function to see what that does.
The `close` function is written in a "defensive" style because it can't be sure what succeeded earlier, so it doesn't make any assumptions.  So if there was no db connection acquired, it returns. However if we got past that it runs a defer to ensures that the connection will be set to nil on return. Next it calls `commitOrRollback`.

OK so just like `close` this function is written in a defensive style, so it starts by checking if the transaction was ever created, and if not it can return.

Next it checks if the requestScope has an error records, if it does that elicits a transaction rollback, so that any 'partial' database changes are reverted. The rollback itself can get an error so we have a similar pattern than we have seen before and it records this as a "system" error.  This will log the error and so forth.

Assuming a happy path we will commit the transaction, and add the usual error handling around this.

OK back to the `close` function. Last but not we release the database connection that the request has been using back to the connection pool, by calling the Release function.

This concludes the request / response cycle.

We have walked through the handler for adding stock. The other handlers work in exactly the same way, so I'll leave it as an exercise for you to take a look at them in your own time.  The remove stock handler has an interesting point of difference in that it published a low stock event when stock levels fall below a threshold, take a look for the `emitLowStockEvent` function to examine how that works.

Because all the handlers follow the same patterns, once you know how one works the others become are easy to understand.  Thats the key to wrirting software thats easy for a team of developers to work on - do the same thing the same way, so that you can dive into different parts of the system and retain familiarity with how things work.


Lets see it in action:

Open a terminal and make sure our service is running by `make bootstrap` then `make run`.  Once the server starts up you will see a a message like "INFO beaker is running". In another terminal run some of the following commands:

- `make test-add` Will add stock
- `make test-remove` Will remove stock
- `make test-get` Will show current stock levels
- `make test-events` Will subscribe to the `event.low_stock` NATS subject and show low stock events when they occur.

## Rats'n'mice

We have covered most but not all of the system. There are unit tests in `internal/api/app_test.go`. These tests run a real NATS server in-process, and connect to a real postgres database, which means that the code being executed is close to production code. These tests avoid mocking out layers by using a real application server and database.

To run the tests run `make test`.  

I've also included a `linter` which checks for common mistakes and errors that the compiler doesn't pick up.  Its really slow to run. I haven't investigated why that is, but thats ok, because I only run it periodically.  Linters are a valuable tool that casts an extra set of eyes over your code.

To make sure my code is always passing, I've included a simple github action in `/workspaces/beaker/.github/workflows/pr-ci.yml` thats triggered whenever we open a pull request. This action checks out the code, build the system, runs the tests and linters.  All these actions happen automatically in the background, and if the tests fail, then my PR will be prevented from being merged.  Its a great way to keep your code in good shape at all times.

# Conclusion.

This concludes this video series on "Production grade system development". We've develoed an API microservice using NATS, Postgres and OpenTelemetry all using the Go language.  We've discussed topics like developement environments, testing, telemetry, performance, and reliability.

I really hope you have got something out of this series and best of luck and blessings to you all in your software development journey. Remember "Iron sharpens iron, and one man sharpens another.”  Thanks for watching
