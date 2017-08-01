package selenium;

//import java.awt.Button;
//import java.io.File;
import java.util.List;
import java.util.concurrent.TimeUnit;

import org.junit.Assert;
import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.chrome.ChromeDriver;
//import org.openqa.selenium.ie.InternetExplorerDriver;

public class libraryTest {
	public static boolean isNull(WebElement button) {
		if(button == null)
		{
			return true;
		}
			return false;
	}
	public static WebElement findButtonByCaption(String caption, WebDriver driver) {
		  final List<WebElement> buttons = 
		    driver.findElements(By.className("v-button"));
		  for (final WebElement button : buttons) {
		    if (button.getText().equals(caption)) {
		      return button;
		    }
		  }
		  return null;
		}
		public static WebElement gofindButtonByLink(String link, WebDriver driver) {
			  WebElement button = null;
		      button = driver.findElement(By.linkText(link));
			      return button;

			}
    public static void main(String[] args) {
        // declaration and instantiation of objects/variables
    	//System.setProperty("webdriver.firefox.marionette","C:\\geckodriver.exe");
		//WebDriver driver = new FirefoxDriver();
		//comment the above 2 lines and uncomment below 2 lines to use Chrome
		System.setProperty("webdriver.chrome.driver","C:\\chromedriver.exe");
		WebDriver driver = new ChromeDriver();
    	//File file = new File("C:/IEDriverServer.exe");
    	//System.setProperty("webdriver.ie.driver", file.getAbsolutePath());
        //WebDriver driver = new InternetExplorerDriver();
        //String baseUrl = "https://springboot-library-haunchless-anticonformist.cfapps.io/";
        String baseUrl = "https://library-go-web-app-gnarled-handle.cfapps.io/";
        driver.manage().window().maximize();
        // launch Fire fox and direct it to the Base URL
        driver.get(baseUrl);
        WebElement checkInButton = gofindButtonByLink("Check In Book", driver);
        if(isNull(checkInButton) != true) {
        	checkInButton.click();
        }
        
         /*
         * compare the actual title of the page with the expected one and print
         * the result as "Passed" or "Failed"
         */
        if (driver.getPageSource().contains("Check In a Book")){
        	if(driver.getPageSource().contains("Date In")) {
        			WebElement submitButton = driver.findElement(By.xpath("//button[contains(.,'Submit')]"));
        		if(isNull(submitButton) != true) {
        			//go back and search for book
        			driver.navigate().back();
        			WebElement searchBar = driver.findElement(By.className("form-control"));
        			if(searchBar.getAttribute("placeholder").equalsIgnoreCase("search")) {
        				searchBar.sendKeys("Book");
        				searchBar.submit();
        				System.out.println("Test Passed!");
        			}
        		}
        	}
        }
        else {
            System.out.println("Test Failed");
        }

        //close IE
        driver.close();
        
    }
}
