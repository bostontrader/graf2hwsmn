# Introduction to graf2hwsmn

The purpose of graf2hwsmn is to provide a means to POST Grafana alerts to the Huawei Cloud Simple Message Notification 
(SMN) service.

Although the original use case was to get Grafana to send alerts via SMS, the SMN service can be configured much more 
broadly than that.

We need graf2hwsmn because although Grafana gives us the ability to POST alerts to a URL, that's not enough to actually 
send messages via the SMN service.  In order to actually do this, we have to jump through a lot of hoops with dealing 
with Access Keys, Secret Keys, signing the request, putting signatures in headers more!.

graf2hwsmn bridges that gap.  With graf2hwsmn you can:

* Create a configuration file that contains your sensitive Huawei credentials.

* Specify the location of said configuration file via command-line argument.

* Launch graf2hwsmn via systemd.

Given these powers, Grafana notification channels can POST a request that will ultimately get sent by the Huawei SMN 
service.

A caveat... graf2hwsmn is not a general-purpose bridge from any random source that might want to POST to Huawei SMN.  
Instead, it is specifically built only for Grafana.  That said, it shouldn't be very difficult to modify or generalize 
to handle other sources.  I leave doing so as an exercise for the reader.

# Security

There are basically two parts to consider when contemplating the security of this:

* Grafana will POST HTTP to graf2hwsmn running somewhere.

* graf2hwsmn will then POST HTTPS to some endpoint in the Huawei cloud.

The link from graf2hwsmn to Huawei cloud is the easy part to secure.  It's done via https and uses their system of 
Access Keys, Secret Keys, and request signing to do it's work.  Good luck hacking that.  That's probably not the weak 
link.

The link from Grafana to graf2hwsmn is the troublesome part.  By default Grafana will POST HTTP to graf2hwsmn.  No 
TLS and no authentication.  If graf2hwsmn is running on the same machine as Grafana then this is a lesser problem.

If you want to run graf2hwsmn on another machine then one obvious method of increasing security would be to use
ssh port forwarding.

Another possible step to increase this security would be to run graf2hwsmn behind nginx configured to use TLS or 
something similar.