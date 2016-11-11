var CreateAppReducer = function(state={}, action) {
    // the logic of update components!!
    if (Object.keys(state).length === 0 && state.constructor === Object) {
        return {inputAppName: ''};
    } else {
        if (action.type === 'UPDATE_APP_NAME') {
            return Object.assign({}, state, {inputAppName: action.appname});
        } else {
            return state;
        }
    }
};

module.exports = CreateAppReducer;
