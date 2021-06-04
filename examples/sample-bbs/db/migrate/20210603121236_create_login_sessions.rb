class CreateLoginSessions < ActiveRecord::Migration[6.1]
  def change
    create_table :login_sessions do |t|
      t.string :refresh_token
      t.datetime :expires_at
      t.datetime :refresh_expires_at
      t.string :user_name
      t.string :user_id

      t.timestamps
    end
  end
end
