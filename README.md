VLOApp server-app
==================

VLO server app acts as a proxy between front-end and google sheets acting as a *database*.

There's no caching implemented as this will be done using either nginx or custom cloudflare rules.

Running it requires google account with access to script projects.

Endpoints
=========

- /lucky-number - [getLuckyNumber](https://github.com/VLO-GDA/gapp-scripts/blob/master/luckynumber.gs#L8)
- /timetable/hours - [getHours](https://github.com/VLO-GDA/gapp-scripts/blob/master/timetable.gs#L16)
- /timetable/group/:group - [getTimetable(group)](https://github.com/VLO-GDA/gapp-scripts/blob/master/timetable.gs#L36)