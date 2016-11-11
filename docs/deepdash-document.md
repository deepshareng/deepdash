Table of Contents
--------
  * [Introduction](#Introduction)
  * [Common Deepdash Response Objects](#Common Deepdash Response Objects)
    * [ErrorStructure](#ErrorStructure)
    * [SimpleSuccessResponse](#SimpleSuccessResponse)
    * [SimpleResponseData](#SimpleResponseData)
    * [AppsGetReturnType](#AppsGetReturnType)
    * [ChannelGetReturnType](#ChannelGetReturnType)
    * [AggregateResult](#AggregateResult)
    * [AggregateCount](#AggregateCount)
    * [AppListInfo](#AppListInfo)
    * [Event](#Event)
    * [Eventlist](#Eventlist)
    * [TokenFormat](#TokenFormat)
    * [UserFormat](#UserFormat)
    * [AppFormat](#AppFormat)
    * [AppItem](#AppItem)
    * [Tag](#Tag)
    * [AppChannel](#AppChannel)
    * [AppSelectedItems](#AppSelectedItems)
  * [Authentication](#Authentication)
  * [Resource API](#Resource API)
    * [App Resource API](#App Resource API)
      * [Get App List](#Get App List)
      * [Get app detailed register Info](#Get app detailed register Info)
      * [Callback service](#Callback service)
      * [Upload an Icon](#Upload an Icon)
      * [Append attribution url to app](#Append attribution url to app)
      * [Create an app](#Create an app )
      * [Modify an app](#Modify an app )
    * [Channel Resource API](#Channel Resource API)
      * [Get all channels brief info of an app](#Get all channels brief info of an app)
      * [Look up channel specific data](#Look up channel specific data)
      * [Create new channel](#Create new channel)
      * [Delete channel](#Delete channel)
      * [Get all channels user defined of an app](#Get all channels user defined of an app)
      * [change channelurl of a channel user defined](#change channelurl of a channel user defined)
      * [Look up User last select event items collection](#Look up User last select event items collection)
      * [Update User last select event items collection](#Update User last select event items collection)
    * [Event Resource API](#Event Resource API)
      * [Get all channels brief info of an app](#Get all channels brief info of an app)
    * [User Resource API](#User Resource API)
      * [Register a user](#Register a user)
      * [Apply a trial to use deepshare](#Apply a trial to use deepshare)
      
## <a name="Introduction"></a>Introduction

This documentation is about the API between deepdash frontend and deepdash backend, and includes three part, authentication part, resource part, view part.



## <a name="Common Deepdash Response Objects"></a>Common Deepdash Response Objects
Deepdash Responses always follow two rules, including two

### <a name="ErrorStructure"></a>ErrorStructure

| Property | Type | Description |
|-----|-----|----|
| `Code` | Int | Error Code, specify the error type |
| `Message` | String | Provide detailed error message |

### <a name="SimpleSuccessResponse"></a> SimpleSuccessResponse
| Property | Type | Description |
|---|---|---|
|`success`|bool||

### <a name="SimpleResponseData"></a> SimpleResponseData
| Property | Type | Description |
|---|---|---|
|`value`|string||

### <a name="AppsGetReturnType"></a> AppsGetReturnType

| Property | Type | Description |
|---|---|---|
| `applist` | Array of [AppListInfo](#AppListInfo) |Array of apps of particular account.|

### <a name="ChannelGetReturnType"></a>ChannelGetReturnType
Return Channel specific data

| Property | Type | Description |
|-----|-----|----|
| `data` | Array of [AggregateResult](#AggregateResult) |  |

### <a name="AggregateResult"></a>AggregateResult

| Property | Type | Description |
|---|---|---|
|`event`|string||
|`counts`|Array of [AggregateCount](#AggregateCount)||

### <a name="AggregateCount"></a>AggregateCount

| Property | Type | Description |
|---|---|---|
|`count`|int||

### <a name="AppListInfo"></a>AppListInfo

| Property | Type | Description |
|-----|-----|----|
| `AppID` |     string  |         |
|`AppName`   |string            ||
|`ChannelInfo` |Array of [AppItem](#AppItem) ||




### <a name="Event"></a>Event
| Property | Type | Description |
|---|---|---|
|`event`|string||
|`display`|string||

### <a name="Eventlist"></a>Eventlist

| Property | Type | Description |
|---|---|---|
|`eventlist`|Array of [Event](#Event)||


> storage formmatting

### <a name="TokenFormat"></a> TokenFormat
| Property | Type | Description |
|---|---|---|
|`passwordid`|string||
|`accountid`|string||
|`token`|string||
|`createat`|time.Time||

### <a name="UserFormat"></a> UserFormat
| Property | Type | Description |
|---|---|---|
|`token`|string||
|`username`|string||
|`password`|string||
|`githubname`|string||
|`realityname`|string||
|`phone`|string||
|`email`|string||
|`wechat`|string||
|`qqaccount`|string||
|`activate`|string||

###<a name="AppFormat"></a> AppFormat


| Property | Type | Description |
|---|---|---|
|`appid`|string||
|`accountid`|string||
|`name`|string||
|`pkg_name`|string||
|`iosbundler`|string||
|`iosscheme`|string||
|`iosdownloadurl`|string||
|`iosunilink`|string||
|`iosteamid`|string||
|`androidpkgname`|string||
|`androidscheme`|string||
|`androidhost`|string||
|`androidapplink`|string||
|`androiddownloadurl`|string||
|`androidisdownloaddirectly`|string||
|`yyburl`|string||
|`yybenable`|string||
|`attributionpushurl`|string||
|`iconurl`|string||
|`download_title`|string||
|`download_msg`|string||

###<a name="AppItem"></a> AppItem

| Property | Type | Description |
|---|---|---|
|`typeid`|string||
|`typename`|string||
|`channelname`|string||
|`userdefine`|string||
|`remark`|string||
|`matchURL`|string||
|`tags`|storage.Tag||
|`appid`|string||

###<a name="Tag"></a>Tag
###<a name="AppChannel"></a>AppChannel

| Property | Type | Description |
|---|---|---|
|`appid`|string||
|`channelname`|string||
|`channelurl`|string||

###<a name="AppSelectedItems"></a>AppSelectedItems
| Property | Type | Description |
|---|---|---|
|`appid`|string||
|`accountid`|string||
|`events`|Array of string||
|`displays`|Array of string||

## <a name="Authentication"></a>Authentication

Authentication uses OAuth, and keep the resource service, and the authentication service in same server.
[OAuth Reference](http://www.ruanyifeng.com/blog/2014/05/oauth_2_0.html)


## <a name="Resource API"></a>Resource API

These Resource APIs always pass the `accountid` attribute through the session. 
In another regarding, the `accountid` in session provides authentication service.


### <a name="App Resource API"></a>App Resource API

Use App resources and related sub-resources API manipulate your apps.

#### <a name="Get App List"></a>Get App List
*handler: GetApplistHandler*		##### Operation	`GET /apps`

Use this call to get a list of app information of particular user.

##### Request 

Pass value specified which account is through session

##### Response 
Return [AppsGetReturnType](#AppsGetReturnType) Object,a list of the brief state of the apps this account owned, including channel possessing, app name, etc.



#### <a name="Get app detailed register Info"></a>Get app detailed register Info
*handler: GetAppinfoHandler*
##### Operation
`GET /apps/:appid`

Use this call to get an app information.

##### Request 

Pass the app id in the URI of a APP call(endpoint)

##### Response 

return a [AppFormat](#AppFormat) object



#### <a name="Callback service"></a>Callback service
*handler: PostCallBackUrlHandler*
##### Operation
`POST /apps/:appid/callback`
##### Request
| Property | Type | Description |
|---|---|---|
|`data`|string||
|`url`|string||

##### Response
Return a [SimpleSuccessResponse](#SimpleSuccessResponse]) Object, success field specifies whether the request is successful or not.

#### <a name="Upload an Icon"></a>Upload an Icon
*handler: UploadIconHandler*
Upload an icon to qiniu, then update the app's icon image, return the icon url.
##### Operation
`POST /action/:appid/uploadicon`
##### Request
Pass appid attribute in the endpoint

| Property | Type | Description |
|---|---|---|
|`uploadfile`|File||
|`url`|string||

##### Response
Return [SimpleResponseData](#SimpleResponseData) Object, value field specifes the icon url.




#### <a name="Append attribution url to app"></a>Append attribution url to app
*handler: PutAppUrlHandler*
##### Operaton
`PUT /apps/:appid/url`
##### Request
Pass appid attribute in the endpoint

| Property | Type | Description |
|---|---|---|   
|`attributionpushurl`|string||
##### Response
Return [SimpleSuccessResponse](#SimpleSuccessResponse), success field demonstrates whether request is ok or not.

#### <a name="Create an app"></a>Create an app 
*handler: PostAppHandler*
##### Operation
`Post /apps`
##### Request
| Property | Type | Description |
|---|---|---|
|`appName`|string||
|`pkgName`|string||
|`iosBundler`|string||
|`iosScheme`|string||
|`iosDownloadUrl`|string||
|`iosUniversalLinkEnable`|string||
|`iosTeamID`|string||
|`iosYYBEnableBelow9`|string||
|`iosYYBEnableAbove9`|string||
|`androidPkgname`|string||
|`androidScheme`|string||
|`androidHost`|string||
|`androidSHA256`|string||
|`androidDownloadUrl`|string||
|`androidIsDownloadDirectly`|string||
|`androidYYBEnable`|string||
|`yyburl`|string||
|`yybenable`|string||
|`attributionPushUrl`|string||
|`iconUrl`|string||
|`downloadTitle`|string||
|`downloadMsg`|string||
|`theme`|string||
|`userConfBgWeChatAndroidTipUrl`|string||
|`userConfBgWeChatIosTipUrl`|string||
##### Response
Return [SimpleResponseData](#SimpleResponseData), value field returns the appid generated. Also set appid in session.






#### <a name="Modify an app"></a>Modify an app 
*handler: PutAppHandler*
##### Operation
`Put /apps/:appid`
##### Request
| Property | Type | Description |
|---|---|---|   
|`appName`|string||
|`pkgName`|string||
|`iosBundler`|string||
|`iosScheme`|string||
|`iosDownloadUrl`|string||
|`iosUniversalLinkEnable`|string||
|`iosTeamID`|string||
|`iosYYBEnableBelow9`|string||
|`iosYYBEnableAbove9`|string||
|`androidPkgname`|string||
|`androidScheme`|string||
|`androidHost`|string||
|`androidSHA256`|string||
|`androidDownloadUrl`|string||
|`androidIsDownloadDirectly`|string||
|`androidYYBEnable`|string||
|`yyburl`|string||
|`yybenable`|string||
|`attributionPushUrl`|string||
|`iconUrl`|string||
|`downloadTitle`|string||
|`downloadMsg`|string||
|`theme`|string||
|`userConfBgWeChatAndroidTipUrl`|string||
|`userConfBgWeChatIosTipUrl`|string||
##### Response
Return [SimpleSuccessResponse](#SimpleSuccessResponse), success field demonstrates whether request is OK or not.





















### <a name="Channel Resource API"></a>Channel Resource API

Use channel resources and related sub-resources API manipulate your channel of apps.

#### <a name="Get all channels brief info of an app"></a>Get all channels brief info of an app
*handler: GetAppChannelsHandler*
#####Operation
`GET /apps/:appid/channels/statistics`
Get all channels of an app, and their simple information, event total value.

#####Request
Values passed in the endpoint: appid

| Property | Type | Description |
|---|---|---| 
|`appid`|string||
|`event`|string|A string with delimiter ',', every field denotes a event we want for this channel, in other words, it is a event filters for all events channel possesses|

#####Response
| Property | Type | Description |
|---|---|---| 
|`typename`|string||
|`channelname`|string||
|`remark`|string||
|`eventPairs`|string|eventPairs could be any value of event name, i.e. match/install_with_params|

#### <a name="Look up channel specific data"></a>Look up channel specific data
*handler: GetChannelStatisticsHandler*
#####Operation
`GET /apps/:appid/statistics`
Significant request, request specific data of a channel. For instance, you can request some period install event of a channel. 
#####Request
Support several format rules requesting data, for pattern, we check every rule in turn.
If there is more than one match, we will return the first one matched.
1. start=x&end=y
2. start=x&limit=y
3. end=x&limit=y
4. limit=x

Values passed in the endpoint: appid, channel

| Property | Type | Description |
|---|---|---| 
|`appid`|string||
|`channel`|string||
|`gran`|string||
|`limit`|string||
|`start`|string||
|`end`|string||
#####Response
Return a [ChannelGetReturnType](#ChannelGetReturnType) Object

#### <a name="Create new channel"></a>Create new channel
*handler: PostNewAppChannelHandler*
#####Operation
`POST /apps/:appid/channels"`
Create a channel for user in deepdash.
#####Request
Values passed in the endpoint: appid

| Property | Type | Description |
|---|---|---| 
|`appid`|string||
|`channelname`|string||
|`channelurl`|string||

#####Response
Return [SimpleSuccessResponse](#SimpleSuccessResponse), success field demonstrates whether request is OK or not.

#### <a name="Delete channel"></a>Delete channel
*handler: DeleteAppChannelHandler*
#####Operation
`DELETE /apps/:appid/channels/:channelname"`
Delete a channel for user in both deepstats/deepdash.
#####Request
Values passed in the endpoint: appid, channelname

| Property | Type | Description |
|---|---|---| 
|`appid`|string||
|`channelname`|string||

#####Response
Return [SimpleSuccessResponse](#SimpleSuccessResponse), success field demonstrates whether request is OK or not.



#### <a name="Get all channels user defined of an app"></a>Get all channels user defined of an app
*handler: GetAppChannelInfoHandler*
#####Operation
`GET /apps/:appid/channels`
#####Request
Values passed in the endpoint: appid

| Property | Type | Description |
|---|---|---| 
|`appid`|string||

#####Response
Return a list of AppChannel object

| Property | Type | Description |
|---|---|---| 
|`data`|Array of [AppChannel](#AppChannel)||





#### <a name="change channelurl of a channel user defined"></a>change channelurl of a channel user defined
*handler: PutChannelUrlHandler*
#####Operation
`PUT /apps/:appid/channelurl`
#####Request
Values passed in the endpoint: appid

| Property | Type | Description |
|---|---|---| 
|`appid`|string||
|`channelname`|string|Specify the particular channel|
|`channelurl`|string|The channel URL you want to change|

#####Response
Return [SimpleSuccessResponse](#SimpleSuccessResponse), success field demonstrates whether request is OK or not.


#### ~~~Get appid types~~~
~~~*handler: TypesGetHandler*~~~
<!---####Operation
`PUT /apps/:appid/channelurl`
#####Request
#####Response
--->

#### <a name="Look up User last select event items collection"></a>Look up User last select event items collection
*handler: GetSelectedItemsHandler*
#####Operation
`GET /apps/:appid/selected_items`
#####Request
Values passed in the endpoint: appid, channel

| Property | Type | Description |
|---|---|---| 
|`appid`|string||
#####Response
Return an [AppSeletedItems](#AppSeletedItems) Object

#### <a name="Update User last select event items collection"></a>Update User last select event items collection
*handler: PutSelectedItemsHandler*
#####Operation
`PUT /apps/:appid/selected_items`
#####Request
An [AppSeletedItems](#AppSeletedItems) Object
#####Response
Return [SimpleSuccessResponse](#SimpleSuccessResponse), success field demonstrates whether request is OK or not.



### <a name="Event Resource API"></a>Event Resource API

This API provides you interface to manipulate your events of apps.

#### <a name="Get all channels brief info of an app"></a>Get all channels brief info of an app
*handler: GetAppEventsHandler*
Get all event of an app
#####Operation
`GET /apps/:appid/events`

#####Request
Values passed in the endpoint: appid, channel

| Property | Type | Description |
|---|---|---| 
|`appid`|string||

#####Response
Return an [Eventlist](#Eventlist) Object


### <a name="User Resource API"></a>User Resource API


#### <a name="Register a user"></a>Register a user
*handler: PostUserHandler*
#####Operation
`POST /this-is-a-clandestine-resource/users`

#####Request
Provide basic user info for registration.

| Property | Type | Description |
|---|---|---| 
|`username`|string||
|`password`|string||

#####Response
Return [SimpleSuccessResponse](#SimpleSuccessResponse), success field demonstrates whether request is OK or not.



#### ~~~Temporary exposed register interface(used by developer only)~~~
~~~*handler: PutUserHandler*~~~
#####~~~Operation~~~
~~~`PUT /this-is-a-clandestine-resource/users`~~~

#### <a name="Apply a trial to use deepshare"></a>Apply a trial to use deepshare
*handler: ApplyTryHandler*
Send mail to administrator applying trial
#####Operation
`POST /apply-try/users`

#####Request
Provide basic app/user info for application.

| Property | Type | Description |
|---|---|---| 
|`appname`|string|App name|
|`appactive`|string|DAU|
|`iosdownloadurl` | string | iOS download URL|
|`androiddownloadurl` |string | android download URL|
|`emailaddress` | string | user email address|
|`phonenumber` | string| user phone number|

#####Response
Return [SimpleSuccessResponse](#SimpleSuccessResponse), success field demonstrates whether request is OK or not.
