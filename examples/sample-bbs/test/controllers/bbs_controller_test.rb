require "test_helper"

class BbsControllerTest < ActionDispatch::IntegrationTest
  test "should get show" do
    get bbs_show_url
    assert_response :success
  end

  test "should get add" do
    get bbs_add_url
    assert_response :success
  end
end
