var AppReducer = function(state={}, action) {
    // the logic of update components!!
    if (Object.keys(state).length === 0 && state.constructor === Object) {
        return {appid: sessionAppId};
    } else {
        if (action.type === 'UPDATE_APP_ID') {
            return Object.assign({}, state, {appid: action.appid});
        } else {
            return state;
        }
    }
};

module.exports = AppReducer;
