var SmsReducer = function(state={}, action) {
    if (Object.keys(state).length === 0 && state.constructor === Object) {
        return {'smses': []};
    } else {
        if (action.type === 'UPDATE_SMS_INFO') {
            return Object.assign({}, state, {'smses': action.smsInfo});
        } else if (action.type === 'UPDATE_SMS_ORDER') {
            var smses = state.smses;
            var sortType = action.sortType;

            var sortDataBy = function(type) {
                switch(type) {
                    case "time" :
                        smses.sort(sortDataByTime);
                        break;
                    default:
                        smses.sort(sortDataByTime);
                        break;
                }

                return Object.assign({}, state, {'smses': smses});
            };

            var sortDataByTime = function(a, b) {
                if (!a.createat) {
                    return -1;
                }
                if (!b.createat) {
                    return 1;
                }
                return parseInt(b.createat, 10) - parseInt(a.createat, 10);
            };

            return sortDataBy(sortType);
        } else {
            return state;
        }
    }
};

module.exports = SmsReducer;
