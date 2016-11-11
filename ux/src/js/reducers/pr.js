var PrReducer = function(state={}, action) {
    if (Object.keys(state).length === 0 && state.constructor === Object) {
        return {'links': []};
    } else {
        if (action.type === 'UPDATE_PR_INFO') {
            return Object.assign({}, state, {'links': action.prInfo});
        } else if (action.type === 'UPDATE_PR_ORDER') {
            var links = state.links;
            var sortType = action.sortType;

            var sortDataBy = function(type) {
                switch(type) {
                    case "time" :
                        links.sort(sortDataByTime);
                        break;
                    case "name" :
                        links.sort(sortDataByName);
                        break;
                    case "type" :
                        links.sort(sortDataByType);
                        break;
                    default:
                        links.sort(sortDataByTime);
                        break;
                }

                return Object.assign({}, state, {'links': links});
            };

            var sortDataByUrl = function(a, b) {
                return a.channelurl.localeCompare(b.channelurl);
            }

            var sortDataByTime = function(a, b) {

                if ((!a.createat) && (!b.createat)) {

                    return sortDataByUrl(a, b);
                }
                if (!a.createat) {
                    return -1;
                }
                if (!b.createat) {
                    return 1;
                }
                if (a.createat === b.createat) {

                    return sortDataByUrl(a, b);
                }
                return parseInt(b.createat, 10) - parseInt(a.createat, 10);
            };


            var sortDataByName = function(a, b) {
                var aRes = a.channelname.split("_");
                var bRes = b.channelname.split("_");
                if (aRes[1] == bRes[1]) {
                    return sortDataByUrl(a, b);
                }
                return aRes[1].localeCompare(bRes[1])
            };

            var sortDataByType = function(a, b) {
                aRes = a.channelname.split("_");
                bRes = b.channelname.split("_");
                if (aRes[0] == bRes[0]) {
                    return sortDataByUrl(a, b)
                }
                return aRes[0].localeCompare(bRes[0])
            }


            return sortDataBy(sortType);
        } else {
            return state;
        }
    }
};

module.exports = PrReducer;
