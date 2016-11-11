var AppInfoReducer = function(state={}, action) {
    // the logic of update components!!
    if (Object.keys(state).length === 0 && state.constructor === Object) {
        return {_: ''};
    } else {
        if (action.type === 'INIT_APP_INFO') {
            return Object.assign({}, state, action.appInfo);
        } else {
            return state;
        }
    }
};

module.exports = AppInfoReducer;
