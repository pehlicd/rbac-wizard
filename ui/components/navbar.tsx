"use client";
import {
	Navbar as NextUINavbar,
	NavbarContent,
	NavbarMenu,
	NavbarMenuToggle,
	NavbarBrand,
	NavbarItem,
	NavbarMenuItem,
} from "@nextui-org/navbar";
import { Link } from "@nextui-org/link";

import { link as linkStyles } from "@nextui-org/theme";

import { siteConfig } from "@/config/site";
import NextLink from "next/link";
import clsx from "clsx";

import { ThemeSwitch } from "@/components/theme-switch";
import {
	GithubIcon,
	Logo,
} from "@/components/icons";

import React from "react";

export function Navbar() {
	const [isMenuOpen, setIsMenuOpen] = React.useState(false);
	return (
		<>
			<NextUINavbar
				isMenuOpen={isMenuOpen} onMenuOpenChange={setIsMenuOpen}
				maxWidth="xl" position="sticky" shouldHideOnScroll isBordered
				classNames={{
					item: [
						"flex",
						"relative",
						"h-full",
						"items-center",
						"data-[active=true]:after:content-['']",
						"data-[active=true]:after:absolute",
						"data-[active=true]:after:bottom-0",
						"data-[active=true]:after:left-0",
						"data-[active=true]:after:right-0",
						"data-[active=true]:after:h-[2px]",
						"data-[active=true]:after:rounded-[2px]",
						"data-[active=true]:after:bg-primary",
					],
				}}
			>
				<NavbarContent className="basis-1/5 sm:basis-full" justify="start">
					<NavbarBrand as="li" className="gap-3 max-w-fit">
						<NextLink className="flex justify-start items-center gap-1" href="/">
							<Logo className="w-12 h-12 m-3" />
							<p className="font-bold text-inherit">RBAC Wizard</p>
						</NextLink>
					</NavbarBrand>
					<ul className="hidden lg:flex gap-4 justify-start ml-2">
						{siteConfig.navItems.map((item) => (
							<NavbarItem key={item.href}>
								<NextLink
									className={clsx(
										linkStyles({ color: "foreground" }),
										"data-[active=true]:text-primary data-[active=true]:font-medium"
									)}
									color="foreground"
									href={item.href}
									onClick={() => setIsMenuOpen(false)}
								>
									{item.label}
								</NextLink>
							</NavbarItem>
						))}
					</ul>
				</NavbarContent>

				<NavbarContent
					className="hidden sm:flex basis-1/5 sm:basis-full"
					justify="end"
				>
					<NavbarItem className="hidden sm:flex gap-2">
						<Link isExternal href={siteConfig.links.github} aria-label="Github">
							<GithubIcon className="text-default-500" />
						</Link>
						<ThemeSwitch />
					</NavbarItem>
				</NavbarContent>

				<NavbarContent className="sm:hidden basis-1 pl-4" justify="end">
					<Link isExternal href={siteConfig.links.github} aria-label="Github">
						<GithubIcon className="text-default-500" />
					</Link>
					<ThemeSwitch />
					<NavbarMenuToggle 
						className="text-default-500" 
						aria-label="Menu"
					/>
				</NavbarContent>

				<NavbarMenu>
					<div className="mx-4 mt-2 flex flex-col gap-2">
						{siteConfig.navMenuItems.map((item, index) => (
							<NavbarMenuItem key={`${item}-${index}`}>
								<Link
									color={
										index === 2
											? "primary"
											: index === siteConfig.navMenuItems.length - 1
											? "danger"
											: "foreground"
									}
									href={item.href}
									size="lg"
									onClick={() => setIsMenuOpen(false)}
									onPress={() => setIsMenuOpen(false)}
								>
									{item.label}
								</Link>
							</NavbarMenuItem>
						))}
					</div>
				</NavbarMenu>
			</NextUINavbar>
		</>
	);
};
