var OverallReducer = function(state={}, action) {
    if (Object.keys(state).length === 0 && state.constructor === Object) {
        return {
            gran: 'd',
            totalPv: 0,
            totalClick: 0,
            totalActive: 0,
            totalInstall: 0,
            oSPlatform: 'all', // all/ios/android
            chartDataByGran:{
                't': {
                    chartXAxis: [],
                    chartDataActive: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataInstall: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataPv: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataClick: {
                        all: [],
                        ios: [],
                        android: [],
                    }
                },
                'd': {
                    chartXAxis: [],
                    chartDataActive: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataInstall: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataPv: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataClick: {
                        all: [],
                        ios: [],
                        android: [],
                    }
                },
                'w': {
                    chartXAxis: [],
                    chartDataActive: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataInstall: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataPv: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataClick: {
                        all: [],
                        ios: [],
                        android: [],
                    }
                },
                'm': {
                    chartXAxis: [],
                    chartDataActive: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataInstall: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataPv: {
                        all: [],
                        ios: [],
                        android: [],
                    },
                    chartDataClick: {
                        all: [],
                        ios: [],
                        android: [],
                    }
                }
            }
            ,
        };
    } else {
        if (action.type === 'UPDATE_OVERALL_INFO') {
            return Object.assign({}, state, action.conf);
        } else {
            return state;
        }
    }
};

module.exports = OverallReducer;
