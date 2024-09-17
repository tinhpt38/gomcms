export const getOsBrowserDevice = (userAgent) => {
    let os = "Unknown OS";
    let browser = "Unknown Browser";
    let device = "Desktop/Laptop"; // Default to desktop, override if mobile

    // Phát hiện hệ điều hành
    if (userAgent.includes("Android")) {
        const androidVersion = userAgent.match(/Android\s([0-9\.]+)/);
        os = androidVersion ? `Android ${androidVersion[1]}` : "Android";
    } else if (userAgent.includes("iPhone")) {
        const iosVersion = userAgent.match(/OS\s([0-9_]+)/);
        os = iosVersion ? `iPhone OS ${iosVersion[1].replace(/_/g, ".")}` : "iPhone OS";
    } else if (userAgent.includes("Linux")) {
        os = "Linux";
    } else if (userAgent.includes("Mac OS X")) {
        os = "Mac OS X";
    } else if (userAgent.includes("Windows")) {
        const windowsVersion = userAgent.match(/Windows NT\s([0-9\.]+)/);
        os = windowsVersion ? `Windows NT ${windowsVersion[1]}` : "Windows";
    }

    // Phát hiện trình duyệt
    if (userAgent.includes("Chrome")) {
        const chromeVersion = userAgent.match(/Chrome\/([0-9\.]+)/);
        browser = chromeVersion ? `Chrome ${chromeVersion[1]}` : "Chrome";
    } else if (userAgent.includes("Safari") && !userAgent.includes("Chrome")) {
        const safariVersion = userAgent.match(/Version\/([0-9\.]+)/);
        browser = safariVersion ? `Safari ${safariVersion[1]}` : "Safari";
    } else if (userAgent.includes("HeyTapBrowser")) {
        const heyTapVersion = userAgent.match(/HeyTapBrowser\/([0-9\.]+)/);
        browser = heyTapVersion ? `HeyTapBrowser ${heyTapVersion[1]}` : "HeyTapBrowser";
    }

    // Phát hiện thiết bị
    if (userAgent.includes("Mobile")) {
        device = "Mobile";
    }

    return {
        os,
        browser,
        device
    };

}

export const formatUserAgent = (userAgent) => {
    const { os, browser, device } = getOsBrowserDevice(userAgent);
    return `${os} - ${browser} - ${device}`;
}
