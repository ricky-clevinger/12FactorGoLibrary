package com.LibraryTest.LibraryTest;
/*
 * Name: libraryTest.java
 * Author: Crystal S. Cox
 * Date: 7/31/2017
 * Description: Test Go Library functionality. Does 'Check In a Book" work as expected? (Pass/Fail)
*/

import java.util.List;
import org.junit.Test;
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
				}//end if
			return false;
		}//end isNull(WebElement)
		public static WebElement findButtonByCaption(String caption, WebDriver driver) {
		  final List<WebElement> buttons = 
		    driver.findElements(By.className("v-button"));
		  for (final WebElement button : buttons) {
		    if (button.getText().equals(caption)) {
		      return button;
		    }//end if
		  }//end for loop
		  return null;
		}//end findButtonByCaption(String, WebDriver)
		public static WebElement gofindButtonByLink(String link, WebDriver driver) {
			  WebElement button = null;
		      button = driver.findElement(By.linkText(link));
			      return button;
		}//end gofindButtonByLink(String, WebDriver)
		@Test
		public void testCheckIn() {
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
	        }//end if
	         /*
	         * check for expected results and print
	         * the results as "Passed" or "Failed"
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
	        			}//end if
	        		}//end if
	        	}//end if
	        }//end if
	        else {
	            System.out.println("Test Failed");
	        }//end else
	        //close IE
	        driver.close();
	}//end testCheckIn()
}//end class