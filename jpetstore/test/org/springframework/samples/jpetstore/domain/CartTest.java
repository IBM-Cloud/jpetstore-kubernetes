package org.springframework.samples.jpetstore.domain;

import static org.junit.Assert.*;
import org.springframework.samples.jpetstore.domain.Cart;
import org.springframework.samples.jpetstore.domain.Item;
import org.springframework.beans.support.PagedListHolder;

import java.util.Iterator;

import org.junit.After;
import org.junit.AfterClass;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;

public class CartTest {

	@BeforeClass
	public static void setUpBeforeClass() throws Exception {
	}

	@AfterClass
	public static void tearDownAfterClass() throws Exception {
	}

	@Before
	public void setUp() throws Exception {
	}

	@After
	public void tearDown() throws Exception {
	}

	@Test
	public void newCart() {
		//fail("Not yet implemented");
		Cart cart = new Cart();
		assertEquals(0, cart.getNumberOfItems());
	}
	
	@Test (expected = NullPointerException.class)
	public void addInvalidItemToCart() {
		Item anItem = null;
		Cart cart = new Cart();
		cart.addItem(anItem,  false);
	}
	
	@Test
	public void addItemToCart() {
		Item anItem = new Item();
		Cart cart = new Cart();
		cart.addItem(anItem,  false);
		assertEquals(1, cart.getNumberOfItems());
	}
	
	@Test
	public void addItemToCartVerify() {
		Item anItem = new Item();
		Cart cart = new Cart();
		cart.addItem(anItem,  false);
		assertEquals(1, cart.getNumberOfItems());
		PagedListHolder list = cart.getCartItemList();
		CartItem rItem = (CartItem)list.getSource().get(0);
		assertEquals(0, rItem.getTotalPrice(), 0.1);
	}

}
