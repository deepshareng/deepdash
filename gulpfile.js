var gulp = require('gulp');
var coffee = require('gulp-coffee');
var browserify = require('gulp-browserify');
var browserify2 = require('browserify');

var runSequence = require('run-sequence');
var minifycss = require('gulp-minify-css');
var plumber = require('gulp-plumber');
var jshint = require('gulp-jshint');
var rename = require('gulp-rename');
var uglify = require('gulp-uglify');
var stylus = require('gulp-stylus');
var jade = require('gulp-jade');
var babel = require('gulp-babel');
var nib = require('nib');
var source = require('vinyl-source-stream');
var buffer = require('vinyl-buffer');
var babelify = require('babelify');

gulp.task('jshint', function() {
    return gulp.src(['ux/src/js/lib/**/*.js', 'ux/src/js/page/**/*.js'])
        .pipe(jshint())
        .pipe(jshint.reporter('default'));
});

gulp.task('browserify', function() {
    return gulp.src(['ux/src/js/page/**/*.js', '!ux/src/js/page/login.js'])
        // 防止代码出错造成watch任务跳出
        .pipe(plumber()) 
        .pipe(browserify())
        //.pipe(babel({presets: ['es2015']}))
        .pipe(gulp.dest('ux/assets/v2/js'))
});

gulp.task('react', function() {
    process.env.NODE_ENV = 'production';
    return browserify2('./ux/src/js/page/index.jsx') 
        .transform(babelify, {presets: ['es2015', 'react']})
        .bundle()
        .pipe(source('index.js'))
        .pipe(buffer())
        .pipe(uglify())
        .pipe(gulp.dest('ux/assets/v2/js'));
});

gulp.task('uglify-login', ['jshint'], function() {
    return gulp.src('ux/src/js/page/login.js')
        .pipe(uglify())
        .pipe(gulp.dest('ux/assets/v2/js'));
});

gulp.task('html', function() {
    return gulp.src('ux/src/jade/*.jade')
        .pipe(plumber())
        .pipe(jade({
            pretty: true
        }))
        .pipe(gulp.dest('ux/html'))
});

gulp.task('css', function() {
    return gulp.src('ux/src/stylus/page/*.styl')
        .pipe(plumber())
        .pipe(stylus({
            use: nib(),
            import: ['nib']
        }))
        .pipe(gulp.dest('ux/assets/v2/css'))
        //.pipe(minifycss({compatibility:'ie7'}))
});

gulp.task('build', function() {
    runSequence('html', 'css', 'browserify');
});

gulp.task('deploy', function() {
    gulp.src('ux/assets/v2/js/*.js')
        .pipe(gulp.dest('ds/assets/v2/js'));
    gulp.src('ux/assets/v2/css/*.css')
        .pipe(gulp.dest('ds/assets/v2/css'));
    gulp.src('ux/assets/v2/fonts/**/*')
        .pipe(gulp.dest('ds/assets/v2/fonts'));
    gulp.src('ux/assets/v2/images/**/*')
        .pipe(gulp.dest('ds/assets/v2/images'));
    gulp.src('ux/html/index.html')
        .pipe(gulp.dest('ds/site/'));
    gulp.src('ux/html/unverified-email.html')
        .pipe(gulp.dest('ds/site/'));
    gulp.src('ux/html/private-add-user.html')
        .pipe(gulp.dest('ds/site/'));
    gulp.src('ux/html/private-permission.html')
        .pipe(gulp.dest('ds/site/'));
});

gulp.task('default', function() {
    runSequence('html', 'css', 'browserify', 'react', 'uglify-login', 'deploy');
});

gulp.task('withouReact', function() {
    runSequence('html', 'css', 'browserify', 'uglify-login', 'deploy');
});
