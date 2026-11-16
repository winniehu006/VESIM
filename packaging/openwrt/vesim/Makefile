include $(TOPDIR)/rules.mk

PKG_NAME:=vesim
VESIM_VERSION?=2.1.23
VESIM_PACKAGE_VERSION?=$(VESIM_VERSION)
PKG_VERSION:=$(VESIM_PACKAGE_VERSION)
PKG_RELEASE:=2

PKG_LICENSE:=SEE LICENSE
PKG_MAINTAINER:=winniehu006 <https://github.com/winniehu006/vesim>

ifeq ($(ARCH),x86_64)
  VESIM_BINARCH:=amd64
else ifeq ($(ARCH),aarch64)
  VESIM_BINARCH:=arm64
else ifeq ($(ARCH),arm)
  VESIM_BINARCH:=armv7
else
ifneq ($(DUMP),1)
  $(error VESIM OpenWrt package supports x86_64, aarch64, and arm)
endif
endif

PKG_SOURCE:=vesim_v$(VESIM_VERSION)_openwrt_$(VESIM_BINARCH)
PKG_SOURCE_URL:=https://github.com/winniehu006/vesim/releases/download/v$(VESIM_VERSION)
# Supplied by build-packages.sh from the verified local release artifact.
PKG_HASH:=$(VESIM_BINARY_SHA256)

include $(INCLUDE_DIR)/package.mk

define Package/vesim
  SECTION:=net
  CATEGORY:=Network
  TITLE:=VESIM modem management
  URL:=https://github.com/winniehu006/vesim
  DEPENDS:=+libqmi +kmod-usb-net-qmi-wwan +kmod-usb-serial-option
endef

define Package/vesim/description
  Prebuilt musl-static VESIM for OpenWrt. Do not use the UPX GitHub
  linux_* binaries; those need glibc.
endef

define Package/vesim/conffiles
/etc/vesim/config.yaml
endef

define Package/vesim-modem-voice
  SECTION:=net
  CATEGORY:=Network
  TITLE:=VESIM optional module direct calling dependencies
  DEPENDS:=+vesim +vesim-adb +alsa-utils +kmod-usb-audio
endef

define Package/vesim-modem-voice/description
  Optional USB audio and modern private ADB for supported module firmware.
  WiFi calling and other VESIM features do not require this package.
endef

define Build/Prepare
	test -n "$(PKG_HASH)"
	printf '%s  %s\n' "$(PKG_HASH)" "$(DL_DIR)/$(PKG_SOURCE)" | sha256sum -c -
	mkdir -p $(PKG_BUILD_DIR)
	cp $(DL_DIR)/$(PKG_SOURCE) $(PKG_BUILD_DIR)/vesim
	chmod 0755 $(PKG_BUILD_DIR)/vesim
endef

define Build/Compile
endef

define Package/vesim/install
	$(INSTALL_DIR) $(1)/usr/bin $(1)/etc/vesim $(1)/etc/init.d $(1)/var/lib/vesim
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/vesim $(1)/usr/bin/vesim
	$(INSTALL_CONF) ./files/config.yaml $(1)/etc/vesim/config.yaml
	$(INSTALL_BIN) ./files/vesim.init $(1)/etc/init.d/vesim
endef

define Package/vesim-modem-voice/install
	$(INSTALL_DIR) $(1)/usr/share/vesim
	$(INSTALL_DATA) ./files/modem-voice.txt $(1)/usr/share/vesim/modem-voice.txt
endef

$(eval $(call BuildPackage,vesim))
$(eval $(call BuildPackage,vesim-modem-voice))
