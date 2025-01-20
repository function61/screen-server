// disable welcome page
pref('browser.aboutwelcome.enabled', false);

// disable nag when entering about:config
pref('browser.aboutConfig.showWarning', false);

// https://stackoverflow.com/a/47353456
pref('datareporting.policy.firstRunURL', '');

pref('browser.ctrlTab.recentlyUsedOrder', false);

// send less data to Mozilla
pref('app.shield.optoutstudies.enabled', false);

// dark mode
// https://superuser.com/a/1694944
pref('layout.css.prefers-color-scheme.content-override', 0);
