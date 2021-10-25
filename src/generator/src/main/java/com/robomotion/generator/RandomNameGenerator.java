package com.robomotion.generator;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

import com.github.javafaker.Faker;
import com.robomotion.app.Context;
import com.robomotion.app.FieldAnnotations;
import com.robomotion.app.Icons;
import com.robomotion.app.Node;
import com.robomotion.app.NodeAnnotations;
import com.robomotion.app.RpcError;
import com.robomotion.app.Runtime.InVariable;
import com.robomotion.app.Runtime.OptVariable;
import com.robomotion.app.Runtime.OutVariable;

	@NodeAnnotations.Name(name = "Robomotion.Generator.RandomNameGenerator")
	@NodeAnnotations.Inputs(inputs = 1)
	@NodeAnnotations.Outputs(outputs = 1)
	@NodeAnnotations.Title(title = "Random Name Generator")
	@NodeAnnotations.Color(color = "#f00")
	@NodeAnnotations.Icon(icon = Icons.mdiCreation)
	public class RandomNameGenerator extends Node {

		@FieldAnnotations.Title(title = "Number of Names")
		@FieldAnnotations.Default(scope = "Custom")
		@FieldAnnotations.MessageScope
		@FieldAnnotations.CustomScope
		public InVariable<String> inNumber;

		@FieldAnnotations.Title(title = "Locale")
		@FieldAnnotations.Default(scope="Custom")
		@FieldAnnotations.CustomScope
		@FieldAnnotations.MessageScope
		@FieldAnnotations.Option
		public OptVariable<String> optLocaleValue;
		

		@FieldAnnotations.Title(title = "Result")
		@FieldAnnotations.Default(scope = "Message", name = "result")
		@FieldAnnotations.MessageOnly
		public OutVariable<Object> outResult;

		@Override
		public void OnCreate() {
		}

		@Override
		public void OnMessage(Context ctx) throws Exception {
			String strNumber = inNumber.Get(ctx);
			if (strNumber == "") {
				throw new RpcError("ErrInvalidArg", "Number of Names can not be empty");
			}

			int num =Integer.parseInt(strNumber);  
			if (num <= 0) {
				throw new RpcError("ErrInvalidArg", "Number must be greater than 0");
			}

			String optLocale = optLocaleValue.Get(ctx);
			Faker faker = new Faker();
			if (optLocale != "") {
				faker = new Faker(new Locale(optLocale));
			}
			List<String> names = new ArrayList<String>();
				for(int i = 0;i < num; i++){
					String fullName = faker.name().firstName() + " " + faker.name().lastName();
					names.add(fullName);
				}			
			outResult.Set(ctx, names.toArray());
		

		}

		@Override
		public void OnClose() {

		}
	}
