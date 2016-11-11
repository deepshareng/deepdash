var React = require('react');
var ReactDOM = require('react-dom');

var Router = require('react-router').Router;
var Route = require('react-router').Route;

var IndexRoute = require('react-router').IndexRoute;
var IndexRedirect = require('react-router').IndexRedirect;

var Link = require('react-router').Link;
var hashHistory = require('react-router').hashHistory;

var ConfigureTpl = require('../components/integrate-conf.jsx').ConfigureTpl;
var CreateAppTpl = require('../components/integrate-create.jsx').CreateAppTpl;

var AppTpl = require('../components/integrate-app.jsx').AppTpl;
var IosTpl = require('../components/integrate-app.jsx').IosTpl;
var AndroidTpl = require('../components/integrate-app.jsx').AndroidTpl;

var WebTpl = require('../components/integrate-web.jsx').WebTpl;
var ApiTpl = require('../components/integrate-web.jsx').ApiTpl;
var JsTpl = require('../components/integrate-web.jsx').JsTpl;

var CompleteTpl = require('../components/integrate-complete.jsx').CompleteTpl;

var PrTpl = require('../containers/pr.jsx').PrTpl;
var SmsTpl = require('../containers/sms.jsx').SmsTpl;
var AppInfoTpl = require('../containers/appinfo.jsx').AppInfoTpl;
var OverallTpl = require('../containers/overall.jsx').OverallTpl;
var MarketingTpl = require('../containers/marketing.jsx').MarketingTpl;
var AppEditTpl = require('../containers/editapp.jsx').AppEditTpl;
var SimulatorTpl = require('../components/integrate-simulator.jsx').SimulatorTpl;

var Provider = require('react-redux').Provider;

var createStore = require('redux').createStore;
var combineReducers = require('redux').combineReducers;

var CreateAppReducer = require('../reducers/integrate-create.js');
var ConfReducer = require('../reducers/integrate-conf.js');
var PrReducer = require('../reducers/pr.js');
var SmsReducer = require('../reducers/sms.js');
var WebReducer = require('../reducers/web.js');
var AppReducer = require('../reducers/app.js');
var AppInfoReducer = require('../reducers/appinfo.js');
var OverallReducer = require('../reducers/overall-info.js');
var MarketingReducer = require('../reducers/marketing-info.js');
var SimulatorReducer = require('../reducers/simulator.js').SimulatorReducer;

var store = createStore(combineReducers({
    createInfo: CreateAppReducer,
    confInfo: ConfReducer,
    prInfo: PrReducer,
    smsInfo: SmsReducer,
    webInfo: WebReducer,
    appIntegrateInfo: AppReducer,
    appInfo: AppInfoReducer,
    overallInfo: OverallReducer,
    marketingInfo: MarketingReducer,
    simulatorInfo: SimulatorReducer,
}));

ReactDOM.render(
    <Provider store={store}>
        <Router history={hashHistory}>
            <Route path="/">
                <IndexRedirect to="overall"/>
                <Route path="create-app" component={CreateAppTpl}/>
                <Route path="configure" component={ConfigureTpl}/>
                <Route path="app" component={AppTpl}>
                    <IndexRedirect to="ios"/>
                    <Route path="ios" component={IosTpl}/>
                    <Route path="android" component={AndroidTpl}/>
                </Route>
                <Route path="web" component={WebTpl}>
                    <IndexRedirect to="api"/>
                    <Route path="api" component={ApiTpl}/>
                    <Route path="js" component={JsTpl}/>
                </Route>
                <Route path="complete" component={CompleteTpl}/>
                <Route path="pr" component={PrTpl}/>
                <Route path="sms" component={SmsTpl}/>
                <Route path="appinfo" component={AppInfoTpl}/>
                <Route path="overall" component={OverallTpl}/>
                <Route path="marketing" component={MarketingTpl}/>
                <Route path="editapp" component={AppEditTpl}/>
                <Route path="simulator" component={SimulatorTpl}/>
            </Route>
        </Router>
    </Provider>,
    document.getElementById('index-app')
);


var DSCOMMON = require('../lib/common.js');
$(function() {
    DSCOMMON.initNav();
    DSCOMMON.renderUserInfo();
    DSCOMMON.initLoginModal();
    //$('#app-list-area').addClass('hide');

    DSCOMMON.getApps(function(data) {
        appNames = data.appNames;
        DSCOMMON.selectApp(data.appIndex, function(){
            DSGLOBAL.initComponent();
        });
    });
});


