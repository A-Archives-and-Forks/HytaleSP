# Introduction

This is a loose guide for setting up a 'offline' hytale server so that players can join using HytaleSP's "Fake Online Mode"
the tl;dr is; you just pass ``--auth-mode=insecure`` to the server.jar instead of ``--auth-mode=authenticated``

if you are trying to join via local multiplayer, or using a "game code" they should just work;
as long as your network supports it, as long as all players are using HytaleSP (or simular) 'offline' launchers;

note: setting this flag will- allow anyone to login as any user anywhere;
as the name suggests this is a bit of a security issue;

id recommend using some sort of authentication plugin ontop of this-
if your doing anything other than making a total "anarchy" server ..

# Obtaining the game server 

Hytale as a game is architected very weirdly;

the Hytale client is C# / .NET NativeAOT;
the Server is Java;
the Launcher is Go;

offically, you're supposed to use the HytaleDownloader tool from hytale.com;
unfortunately  this requires authentication throughout the entire thing ..

thankfully though the hytale client also includes a copy of the server, so you can just take it from there,

.. it will be located at : 
windows: ``%APPDATA%\hytLauncher\{patchline}\{verison}\Server`` 
linux: ``~\.config\hytLauncher\{patchline}\{verison}\Server`` 
 ** you will also need the "Assets.zip" file.

  (these are default locations if you changed where hytaleSP stores files then it could be in differnet places.)

if you prefer to have a more "HytaleDownloader"-like experience; 
 .. can also copy HytaleSP itself to your VPS/Server and download the files with the command line like so : 

 
```
./HytaleSP --patchline=pre-release --version=latest --download-server=./HytaleServer
```

# Obtaining java

if you already have java, then you can mostly skip this step as long as its Java 25 or higher;
you can download java from the [offical website](https://www.oracle.com/java/technologies/downloads/), or using your package manager;

however if you prefer you can also use the version distributed by Hytale themselves;
which can also be downloaded using HytaleSP- using:
```
./HytaleSP --download-jre=./jre
```

# Running the server 

you will need the Java Virtual Machine (JVM) to run the server; 
 the command recommended to run the server with is as follows: 

``java -XX:AOTCache=HytaleServer.aot -jar HytaleServer.jar --assets ../Assets.zip --backup --backup-dir backups --backup-frequency 30 --auth-mode=insecure``

breaking this command down, it is as follows: 

``java`` - the java virtual machine, an interpreter to run java applications

``-XX:AOTCache=HytaleServer.aot`` - ahead-of-time compiled java bytecode, this is precompiling a bunch of code used, which makes the server run a bit faster,

``-jar HytaleServer.jar`` - this specifies the java binary file to be run by the JVM

``--assets ../Assets.zip`` - this specifies the path to the  "assets.zip" file, which is shared by the client and the server ..

``--backup`` - this enables backups 

``--backup-frequency 30`` - this makes a backup every 30 minutes

``--auth-mode=insecure`` - this disables authentication so players can play using the "fakeonline" launch mode.

you can optionally add ``-Xms2G`` or ``-Xmx4G`` to the server if you would like to provide it with more RAM; 
where `Xms` is the minimum ram the server can use, and `Xmx` is the max amount of ram the server can use;
and ``2G`` and ``4G`` are the amounts of ram respectively;


